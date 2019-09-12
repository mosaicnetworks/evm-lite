package state

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	ethState "github.com/ethereum/go-ethereum/core/state"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/sirupsen/logrus"

	bcommon "github.com/mosaicnetworks/evm-lite/src/common"
	"github.com/mosaicnetworks/evm-lite/src/currency"
)

var (
	fdLimit        = 8192
	gasLimit       = uint64(1000000000000000000)
	receiptsPrefix = []byte("receipts-")
)

/*
State is the main object that manages the application-state of evm-lite. It is
used by the Service for read-only operations, and by the Consensus system to
apply new transactions. It manages 3 copies of the underlying datastore:

1) it's own ethState, which is the "official" state, that cannot be arbitrarily
   reverted.
2) the write-ahead-state (was), where the consensus system applies transactions
   before committing them to the main state.
3) the transaction-pool's state, where the Service verifies transactions before
   submitting them to the consensus system.

*/
type State struct {
	db ethdb.Database

	ethState    *ethState.StateDB
	signer      ethTypes.Signer
	chainConfig params.ChainConfig //vm.env is still tightly coupled with chainConfig
	vmConfig    vm.Config
	gasLimit    uint64

	was    *WriteAheadState
	txPool *TxPool

	genesisFile string

	logger *logrus.Entry
}

// NewState creates and initializes a new State object. It reads the genesis
// file to create the initial accounts, including the POA smart-contract.
func NewState(dbFile string, dbCache int, genesisFile string, logger *logrus.Entry) (*State, error) {

	db, err := ethdb.NewLDBDatabase(dbFile, dbCache, fdLimit)
	if err != nil {
		return nil, err
	}

	s := &State{
		db:          db,
		signer:      ethTypes.NewEIP155Signer(CustomChainConfig.ChainID),
		chainConfig: CustomChainConfig,
		vmConfig:    vm.Config{Tracer: vm.NewStructLogger(nil)},
		genesisFile: genesisFile,
		logger:      logger,
	}

	if err := s.InitState(); err != nil {
		return nil, err
	}

	return s, nil
}

/******************************************************************************/

// InitState initializes the statedb object, the write-ahead state, the
// transaction-pool, and creates genesis accounts.
func (s *State) InitState() error {

	s.gasLimit = gasLimit

	initState := common.Hash{}

	var err error

	s.ethState, err = ethState.New(initState, ethState.NewDatabase(s.db))
	if err != nil {
		return err
	}

	s.was, err = NewWriteAheadState(s.db,
		initState,
		s.signer,
		s.chainConfig,
		s.vmConfig,
		gasLimit,
		s.logger)

	if err != nil {
		return err
	}

	s.txPool = NewTxPool(s.ethState.Copy(),
		s.signer,
		s.chainConfig,
		s.vmConfig,
		s.gasLimit,
		s.logger)

	// Initialize genesis accounts with balance, code, and state
	err = s.CreateGenesisAccounts()
	if err != nil {
		return err
	}

	return err
}

// CreateGenesisAccounts reads the genesis.json file and creates the regular
// pre-funded accounts, as well as the POA smart-contract account.
func (s *State) CreateGenesisAccounts() error {

	genesis, err := s.GetGenesis()
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	// Regular pre-funded accounts
	for addr, account := range genesis.Alloc {
		address := common.HexToAddress(addr)
		if s.empty(address) {
			s.was.ethState.AddBalance(address, math.MustParseBig256(currency.ExpandCurrencyString(account.Balance)))
			s.was.ethState.SetCode(address, common.Hex2Bytes(account.Code))
			for key, value := range account.Storage {
				s.was.ethState.SetState(address, common.HexToHash(key), common.HexToHash(value))
			}
			s.logger.WithField("address", addr).Debug("Adding account")
		}
	}

	// POA smart-contract account
	if string(genesis.Poa.Address) != "" {
		address := common.HexToAddress(genesis.Poa.Address)
		if s.empty(address) {
			s.was.ethState.AddBalance(address, math.MustParseBig256(currency.ExpandCurrencyString(genesis.Poa.Balance)))
			s.was.ethState.SetCode(address, common.Hex2Bytes(genesis.Poa.Code))
			setPOAADDR(genesis.Poa.Address)
			setPOAABI(genesis.Poa.Abi)
			s.logger.WithField("address", genesis.Poa.Address).Debug("Adding POA smart-contract account")
		}
	}

	if _, err = s.Commit(); err != nil {
		return err
	}

	return nil

}

// empty reports whether the account is non-existant or empty
func (s *State) empty(addr common.Address) bool {
	res := s.ethState.Empty(addr)
	s.logger.Debugf("%s Empty? %v", addr.Hex(), res)
	return res
}

/*******************************************************************************
Methods called by Consensus
*******************************************************************************/

// ApplyTransaction decodes a transaction and applies it to the WAS. It is meant
// to be called by the consensus system to apply transactions sequentially.
func (s *State) ApplyTransaction(
	txBytes []byte,
	txIndex int,
	blockHash common.Hash,
	coinbase common.Address) error {

	var t ethTypes.Transaction
	if err := rlp.Decode(bytes.NewReader(txBytes), &t); err != nil {
		s.logger.WithError(err).Error("Decoding Transaction")
		return err
	}
	s.logger.WithField("hash", t.Hash().Hex()).Debug("Decoded tx")

	return s.was.ApplyTransaction(t, txIndex, blockHash, coinbase)
}

// Commit persists all pending state changes (in the WAS) to the DB, and resets
// the WAS and TxPool
func (s *State) Commit() (common.Hash, error) {
	// commit all state changes to the database
	root, err := s.was.Commit()
	if err != nil {
		s.logger.WithError(err).Error("Committing WAS")
		return root, err
	}

	// Reset main ethState
	if err := s.ethState.Reset(root); err != nil {
		s.logger.WithError(err).Error("Resetting main StateDB")
		return root, err
	}
	s.logger.WithField("root", root.Hex()).Debug("Committed")

	// Reset WAS
	if err := s.was.Reset(root); err != nil {
		s.logger.WithError(err).Error("Resetting WAS")
		return root, err
	}
	s.logger.Debug("Reset WAS")

	// Reset TxPool
	if err := s.txPool.Reset(root); err != nil {
		s.logger.WithError(err).Error("Resetting TxPool")
		return root, err
	}
	s.logger.Debug("Reset TxPool")

	return root, nil
}

/*******************************************************************************
Service on Main
*******************************************************************************/

// GetBalance returns an account's balance from the main ethState
func (s *State) GetBalance(addr common.Address) *big.Int {
	return s.ethState.GetBalance(addr)
}

// GetNonce returns an account's nonce from the main ethState
func (s *State) GetNonce(addr common.Address) uint64 {
	return s.ethState.GetNonce(addr)
}

// GetCode returns an account's bytecode from the main ethState
func (s *State) GetCode(addr common.Address) []byte {
	return s.ethState.GetCode(addr)
}

// GetTransaction fetches transactions by hash directly from the DB.
func (s *State) GetTransaction(hash common.Hash) (*ethTypes.Transaction, error) {
	// Retrieve the transaction itself from the database
	data, err := s.db.Get(hash.Bytes())
	if err != nil {
		s.logger.WithError(err).Error("GetTransaction")
		return nil, err
	}
	var tx ethTypes.Transaction
	if err := rlp.DecodeBytes(data, &tx); err != nil {
		s.logger.WithError(err).Error("Decoding Transaction")
		return nil, err
	}

	return &tx, nil
}

// GetReceipt fetches transaction receipts by transaction hash directly from the
// DB
func (s *State) GetReceipt(txHash common.Hash) (*ethTypes.Receipt, error) {
	data, err := s.db.Get(append(receiptsPrefix, txHash.Bytes()...))
	if err != nil {
		s.logger.WithError(err).Error("GetReceipt")
		return nil, err
	}

	var receipt ethTypes.ReceiptForStorage
	if err := rlp.DecodeBytes(data, &receipt); err != nil {
		s.logger.WithError(err).Error("Decoding Receipt")
		return nil, err
	}

	return (*ethTypes.Receipt)(&receipt), nil
}

// GetGasLimit returns the gas limit set between commit calls
func (s *State) GetGasLimit() uint64 {
	return s.gasLimit
}

// GetGenesis reads and unmarshals the genesis.json file
func (s *State) GetGenesis() (bcommon.Genesis, error) {
	if _, err := os.Stat(s.genesisFile); err != nil {
		return bcommon.Genesis{}, err
	}

	contents, err := ioutil.ReadFile(s.genesisFile)
	if err != nil {
		return bcommon.Genesis{}, err
	}

	var genesis bcommon.Genesis

	if err := json.Unmarshal(contents, &genesis); err != nil {
		return bcommon.Genesis{}, err
	}

	return genesis, nil
}

// GetAuthorisingAccount returns the address of the smart contract which handles
// the list of authorized peers
func (s *State) GetAuthorisingAccount() string {
	return POAADDR.String()
}

// GetAuthorisingABI returns the abi of the smart contract which handles the
// list of authorized peers
func (s *State) GetAuthorisingABI() string {
	return POAABISTRING
}

// GetSigner returns the state's signer
func (s *State) GetSigner() ethTypes.Signer {
	return s.signer
}

/*******************************************************************************
Service on WAS
*******************************************************************************/

// Call executes a readonly transaction on a copy of the WAS. It is called
// by the service handlers
func (s *State) Call(callMsg ethTypes.Message) ([]byte, error) {
	s.logger.Debug("Call")

	context := NewContext(callMsg.From(), common.Address{}, 0, big.NewInt(0))

	// We use a copy of the ethState because even call transactions increment
	// the sender's nonce
	vmenv := vm.NewEVM(context, s.was.ethState.Copy(), &s.chainConfig, s.vmConfig)

	res, _, _, err := core.ApplyMessage(vmenv, callMsg, new(core.GasPool).AddGas(gasLimit))
	if err != nil {
		s.logger.WithError(err).Error("Executing Call on WAS")
		return nil, err
	}

	return res, err
}

// CreateReceiptPromise crates a new receipt promise
func (s *State) CreateReceiptPromise(hash common.Hash) *ReceiptPromise {
	p := NewReceiptPromise(hash)

	s.was.receiptPromises[hash] = p

	return p
}

/*******************************************************************************
Service on TxPool
*******************************************************************************/

// CheckTx attempts to apply a transaction to the TxPool's statedb. It is called
// by the Service handlers to check if a transaction is valid before submitting
// it to the consensus system. This also updates the sender's Nonce in the
// TxPool's statedb.
func (s *State) CheckTx(tx *ethTypes.Transaction) error {
	return s.txPool.CheckTx(tx)
}

// GetPoolBalance returns an account's balance from the txpool's ethState
func (s *State) GetPoolBalance(addr common.Address) *big.Int {
	return s.txPool.ethState.GetBalance(addr)
}

// GetPoolNonce returns an account's nonce from the txpool's ethState
func (s *State) GetPoolNonce(addr common.Address) uint64 {
	return s.txPool.ethState.GetNonce(addr)
}

// GetPoolCode returns an account's bytecode from the txpool;s ethState
func (s *State) GetPoolCode(addr common.Address) []byte {
	return s.txPool.ethState.GetCode(addr)
}

/*******************************************************************************
POA
*******************************************************************************/

// CheckAuthorised queries the POA smart-contract to check if the address is
// authorised. It is called by the consensus system when deciding to add or
// remove a peer.
func (s *State) CheckAuthorised(addr common.Address) (bool, error) {

	callData, err := POAABI.Pack("checkAuthorised", addr)
	if err != nil {
		s.logger.Warningf("couldn't pack arguments: %v", err)
	}

	// Apply an ethereum call message (no state update) to query the
	// smart-contract. Since there's no nonce check and the gas price is set to
	// zero, an arbitrary address can be used as the source of the tx.
	ethMsg := ethTypes.NewMessage(POAADDR,
		&POAADDR,
		uint64(1),
		big.NewInt(0),
		s.GetGasLimit(),
		big.NewInt(0),
		callData,
		false)

	s.logger.WithFields(logrus.Fields{
		"addr":     addr.Hex(),
		"callData": hex.EncodeToString(callData),
		"contract": POAADDR.String(),
	}).Debug("checkAuthorised")

	res, err := s.Call(ethMsg)
	if err != nil {
		return false, err
	}

	unpackRes := new(bool)
	POAABI.Unpack(&unpackRes, "checkAuthorised", res)

	if *unpackRes {
		return true, nil
	}

	return false, nil

}
