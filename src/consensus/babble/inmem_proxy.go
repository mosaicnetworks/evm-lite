package babble

import (
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mosaicnetworks/babble/src/hashgraph"
	"github.com/mosaicnetworks/babble/src/proxy"
	"github.com/mosaicnetworks/evm-lite/src/service"
	"github.com/mosaicnetworks/evm-lite/src/state"
	"github.com/sirupsen/logrus"
)

// InmemProxy implements the Babble AppProxy interface
type InmemProxy struct {
	service  *service.Service
	state    *state.State
	submitCh chan []byte
	logger   *logrus.Entry
}

// NewInmemProxy initializes and return a new InmemProxy
func NewInmemProxy(state *state.State,
	service *service.Service,
	submitCh chan []byte,
	logger *logrus.Logger) *InmemProxy {

	return &InmemProxy{
		service:  service,
		state:    state,
		submitCh: submitCh,
		logger:   logger.WithField("module", "babble/proxy"),
	}
}

/*******************************************************************************
Implement Babble AppProxy Interface
*******************************************************************************/

// SubmitCh is the channel through which the Service sends transactions to the
// node.
func (p *InmemProxy) SubmitCh() chan []byte {
	return p.submitCh
}

// CommitBlock commits Block to the State and expects the resulting state hash
func (p *InmemProxy) CommitBlock(block hashgraph.Block) (proxy.CommitResponse, error) {
	p.logger.Debug("CommitBlock")

	blockHashBytes, err := block.Hash()
	blockHash := common.BytesToHash(blockHashBytes)

	for i, tx := range block.Transactions() {
		if err := p.state.ApplyTransaction(tx, i, blockHash); err != nil {
			return proxy.CommitResponse{}, err
		}
	}

	hash, err := p.state.Commit()
	if err != nil {
		return proxy.CommitResponse{}, err
	}

	// Process internal transactions
	// For every join request check whether the peer's public key has been authorised.
	// Apply an ethereum call message (no state update) to query the smart contract.
	// Since there's no nonce check and the gas price is set to zero, an arbitrary address
	// can be used as the source of the tx.

	objABI, _ := abi.JSON(strings.NewReader("[{\"type\":\"function\",\"inputs\": [{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"checkAuthorised\",\"outputs\": [{\"name\":\"\",\"type\":\"bool\"}]}]"))
	fromAddress := common.HexToAddress("0x1337133713371337133713371337133713371337")
	contractAddress := common.HexToAddress(p.state.GetAuthorisingAccount())
	nonce := uint64(1) // not used
	amount := big.NewInt(0)
	gasLimit := p.state.GetGasLimit()
	gasPrice := big.NewInt(0)
	checkNonce := false

	receipts := []hashgraph.InternalTransactionReceipt{}

	for _, tx := range block.InternalTransactions() {

		if tx.Body.Type == hashgraph.PEER_ADD {

			pk, err := crypto.UnmarshalPubkey(tx.Body.Peer.PubKeyBytes())
			if err != nil {
				p.logger.Warningf("couldn't unmarshal pubkey bytes: %v", err)
			}

			addr := crypto.PubkeyToAddress(*pk)

			// var param [32]byte
			// copy(param[12:], addr.Bytes())

			callData, err := objABI.Pack("checkAuthorised", addr)
			if err != nil {
				p.logger.Warningf("couldn't pack arguments: %v", err)
			}

			//0x1a3e994500000000000000000000000038842a05b3dfd507bf14b8c33ea90f747f0e0ec2

			ethMsg := ethTypes.NewMessage(fromAddress,
				&contractAddress,
				nonce,
				amount,
				gasLimit,
				gasPrice,
				callData,
				checkNonce)

			p.logger.WithFields(logrus.Fields{
				"addr":     addr.Hex(),
				"callData": hex.EncodeToString(callData),
				"contract": contractAddress.Hex(),
			}).Debug("checkAuthorised")

			if res, err := p.state.Call(ethMsg); err != nil {
				receipts = append(receipts, tx.AsRefused())
			} else {
				unpackRes := new(bool)
				objABI.Unpack(&unpackRes, "checkAuthorised", res)

				if *unpackRes {
					p.logger.Debug("Accepted peer")
					receipts = append(receipts, tx.AsAccepted())
				} else {
					p.logger.Error("Rejected peer")
					receipts = append(receipts, tx.AsRefused())
				}
			}

		}
	}

	res := proxy.CommitResponse{
		StateHash:                   hash.Bytes(),
		InternalTransactionReceipts: receipts,
	}

	return res, nil
}

//TODO - Implement these two functions
func (p *InmemProxy) GetSnapshot(blockIndex int) ([]byte, error) {
	return []byte{}, nil
}

func (p *InmemProxy) Restore(snapshot []byte) error {
	return nil
}
