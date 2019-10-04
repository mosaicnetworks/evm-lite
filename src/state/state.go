package state

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/sirupsen/logrus"

	bcommon "github.com/mosaicnetworks/evm-lite/src/common"
)

var (
	_fdLimit  = 8192
	_gasLimit = uint64(1000000000000000000)
)

/*
State is the main THREAD SAFE object that manages the application-state of
evm-lite. It is used by the Service for read-only operations, and by the
Consensus system to apply new transactions. It manages 3 copies of the
underlying datastore:

1) it's own state, which is the "official" state, that cannot be arbitrarily
   reverted.
2) the write-ahead-state (was), where the consensus system applies transactions
   before committing them to the main state.
3) the transaction-pool's state, where the Service verifies transactions before
   submitting them to the consensus system.
*/
type State struct {
	main   BaseState
	was    *WriteAheadState
	txPool *TxPool

	genesisFile string

	logger *logrus.Entry
}

// NewState creates and initializes a new State object. It reads the genesis
// file to create the initial accounts, including the POA smart-contract.
func NewState(dbFile string, dbCache int, genesisFile string, logger *logrus.Entry) (*State, error) {

	// db is THREAD SAFE and reused by base, was, and txpool
	db, err := ethdb.NewLDBDatabase(dbFile, dbCache, _fdLimit)
	if err != nil {
		return nil, err
	}

	main := NewBaseState(db,
		common.Hash{},
		ethTypes.NewEIP155Signer(CustomChainConfig.ChainID),
		CustomChainConfig,
		vm.Config{Tracer: vm.NewStructLogger(nil)},
		_gasLimit,
	)

	s := &State{
		main:        main,
		was:         NewWriteAheadState(main.Copy(), logger),
		txPool:      NewTxPool(main.Copy(), logger),
		genesisFile: genesisFile,
		logger:      logger,
	}

	// Initialize genesis accounts with balance, code, and state
	err = s.CreateGenesisAccounts()
	if err != nil {
		return nil, err
	}

	return s, nil
}

/******************************************************************************/

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

		s.was.CreateAccount(address,
			account.Code,
			account.Storage,
			account.Balance)

		s.logger.WithField("address", addr).Debug("Adding account")
	}

	// POA smart-contract account
	if string(genesis.Poa.Address) != "" {
		address := common.HexToAddress(genesis.Poa.Address)

		s.was.CreateAccount(address,
			genesis.Poa.Code,
			map[string]string{},
			genesis.Poa.Balance)

		setPOAADDR(genesis.Poa.Address)
		setPOAABI(genesis.Poa.Abi)

		s.logger.WithField("address", genesis.Poa.Address).Debug("Adding POA smart-contract account")

	}

	if _, err = s.Commit(); err != nil {
		return err
	}

	return nil

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

	t, err := NewEVMLTransaction(txBytes, s.GetSigner())
	if err != nil {
		s.logger.WithError(err).Error("Decoding Transaction")
		return err
	}

	if s.logger.Level > logrus.InfoLevel {
		s.logger.WithField("hash", t.Hash().Hex()).Debug("Decoded tx")
	}

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

	// Reset Main
	if err := s.main.Reset(root); err != nil {
		s.logger.WithError(err).Error("Resetting main StateDB")
		return root, err
	}
	if s.logger.Level > logrus.InfoLevel {
		s.logger.WithField("root", root.Hex()).Debug("Committed")
	}

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
Config
*******************************************************************************/

// GetGasLimit returns the gas limit set between commit calls
func (s *State) GetGasLimit() uint64 {
	return s.main.gasLimit
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
	return s.main.signer
}

/*******************************************************************************
WAS & TxPool
*******************************************************************************/

// Call executes a readonly transaction on a copy of the WAS. It is called by
// the service handlers
func (s *State) Call(callMsg ethTypes.Message) ([]byte, error) {
	res, err := s.was.Call(callMsg)
	if err != nil {
		s.logger.WithError(err).Error("Executing Call on WAS")
		return nil, err
	}

	return res, err
}

// CreateReceiptPromise crates a new receipt promise
func (s *State) CreateReceiptPromise(hash common.Hash) *ReceiptPromise {
	return s.was.CreateReceiptPromise(hash)
}

// CheckTx attempts to apply a transaction to the TxPool's stateDB. It is called
// by the Service handlers to check if a transaction is valid before submitting
// it to the consensus system. This also updates the sender's Nonce in the
// TxPool's statedb.
func (s *State) CheckTx(tx *EVMLTransaction) error {
	return s.txPool.CheckTx(tx)
}

// GetBalance returns an account's balance
func (s *State) GetBalance(addr common.Address, fromPool bool) *big.Int {
	if fromPool {
		return s.txPool.GetBalance(addr)
	}
	return s.main.GetBalance(addr)
}

// GetNonce returns an account's nonce
func (s *State) GetNonce(addr common.Address, fromPool bool) uint64 {
	if fromPool {
		return s.txPool.GetNonce(addr)
	}
	return s.main.GetNonce(addr)
}

// GetCode returns an account's bytecode
func (s *State) GetCode(addr common.Address, fromPool bool) []byte {
	if fromPool {
		return s.txPool.GetCode(addr)
	}
	return s.main.GetCode(addr)
}

// GetTransaction fetches a transaction from the WAS
func (s *State) GetTransaction(txHash common.Hash) (*ethTypes.Transaction, error) {
	return s.was.GetTransaction(txHash)
}

// GetReceipt fetches a transaction's receipt from the WAS
func (s *State) GetReceipt(txHash common.Hash) (*ethTypes.Receipt, error) {
	return s.was.GetReceipt(txHash)
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

	if s.logger.Level > logrus.InfoLevel {
		s.logger.WithFields(logrus.Fields{
			"addr":     addr.Hex(),
			"callData": hex.EncodeToString(callData),
			"contract": POAADDR.String(),
		}).Debug("checkAuthorised")
	}

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
