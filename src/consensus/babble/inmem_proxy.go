package babble

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
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

	internalTransactions := block.InternalTransactions()

	objABI, _ := abi.JSON(strings.NewReader("[{\"type\":\"function\",\"inputs\": [{\"name\":\"pubKey\",\"type\":\"bytes32\"}],\"name\":\"checkAuthorisedPublicKey\",\"outputs\": [{\"name\":\"\",\"type\":\"bool\"}]}]"))
	fromAddress := common.HexToAddress("0x1337133713371337133713371337133713371337") // there's no state update and no nonce check so doesn't matter what address we use
	contractAddress := common.HexToAddress("0xabbaabbaabbaabbaabbaabbaabbaabbaabbaabba")
	nonce := uint64(1) // checkNonce set to false so doesn't matter
	amount := big.NewInt(0)
	gasLimit := uint64(0) // gasLimit of zero should be OK since the gasPrice is zero
	gasPrice := big.NewInt(0)
	checkNonce := false

	for i, tx := range internalTransactions {

		if tx.Type == hashgraph.PEER_ADD {

			callData, err := objABI.Pack("checkAuthorisedPublicKey", []byte(tx.Peer.PubKeyHex)) // check for errors

			if err != nil {
				return proxy.CommitResponse{}, err
			}

			ethMsg := ethTypes.NewMessage(fromAddress,
				&contractAddress,
				nonce,
				amount,
				gasLimit,
				gasPrice,
				callData,
				checkNonce)

			if res, err := p.state.Call(ethMsg); err != nil {
				if err != nil {
					var unpackRes bool
					objABI.Unpack(unpackRes, "checkAuthorisedPublicKey", res)

					if unpackRes {
						internalTransactions[i].Accept()
					} else {
						internalTransactions[i].Refuse()
					}
				} else {
					internalTransactions[i].Refuse()
				}
			}

		}
	}

	res := proxy.CommitResponse{
		StateHash:            hash.Bytes(),
		InternalTransactions: internalTransactions,
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
