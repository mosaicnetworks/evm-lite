package state

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"
	"syscall"

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
)

var (
	gasLimit       = uint64(1000000000000000000)
	txMetaSuffix   = []byte{0x01}
	receiptsPrefix = []byte("receipts-")
	MIPMapLevels   = []uint64{1000000, 500000, 100000, 50000, 1000}
)

type State struct {
	db       ethdb.Database
	ethState *ethState.StateDB
	was      *WriteAheadState
	txPool   *TxPool

	gasLimit           uint64
	authorisingAccount string
	authorisingAbi     string

	signer      ethTypes.Signer
	chainConfig params.ChainConfig //vm.env is still tightly coupled with chainConfig
	vmConfig    vm.Config

	genesisFile string

	logger *logrus.Logger
}

func NewState(logger *logrus.Logger, dbFile string, dbCache int, genesisFile string) (*State, error) {
	handles, err := getFdLimit()
	if err != nil {
		return nil, err
	}

	db, err := ethdb.NewLDBDatabase(dbFile, dbCache, handles)
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

//------------------------------------------------------------------------------

//InitState initializes the statedb object, the write-ahead state, the
//transaction-pool, and creates genesis accounts.
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

	//Initialize genesis accounts with balance, code, and state
	err = s.CreateGenesisAccounts()
	if err != nil {
		return err
	}

	return err
}

//Commit persists all pending state changes (in the WAS) to the DB, and resets
//the WAS and TxPool
func (s *State) Commit() (common.Hash, error) {
	//commit all state changes to the database
	root, err := s.was.Commit()
	if err != nil {
		s.logger.WithError(err).Error("Committing WAS")
		return root, err
	}

	//Reset main ethState
	if err := s.ethState.Reset(root); err != nil {
		s.logger.WithError(err).Error("Resetting main StateDB")
		return root, err
	}
	s.logger.WithField("root", root.Hex()).Debug("Committed")

	//Reset WAS
	if err := s.was.Reset(root); err != nil {
		s.logger.WithError(err).Error("Resetting WAS")
		return root, err
	}
	s.logger.Debug("Reset WAS")

	//Reset TxPool
	if err := s.txPool.Reset(root); err != nil {
		s.logger.WithError(err).Error("Resetting TxPool")
		return root, err
	}
	s.logger.Debug("Reset TxPool")

	return root, nil
}

//------------------------------------------------------------------------------

//Call executes a readonly transaction on the statedb. It is called by the
//service handlers
func (s *State) Call(callMsg ethTypes.Message) ([]byte, error) {
	s.logger.Debug("Call")

	context := NewContext(callMsg.From(), 0, callMsg.GasPrice())

	//We use a copy of the ethState because even call transactions increment the
	//sender's nonce
	vmenv := vm.NewEVM(context, s.was.ethState.Copy(), &s.chainConfig, s.vmConfig)

	// Apply the transaction to the current state (included in the env)
	res, _, _, err := core.ApplyMessage(vmenv, callMsg, new(core.GasPool).AddGas(gasLimit))
	if err != nil {
		s.logger.WithError(err).Error("Executing Call on WAS")
		return nil, err
	}

	return res, err
}

//CheckTx attempt to apply a transaction to the TxPool's statedb. It is called
//by the Service handlers to check if a transaction is valid before submitting
//it to the consensus system. This also updates the sender's Nonce in the
//TxPool's statedb.
func (s *State) CheckTx(tx *ethTypes.Transaction) error {
	return s.txPool.CheckTx(tx)
}

//ApplyTransaction decodes a transaction and applies it to the WAS. It is meant
//to be called by the consensus system to apply transactions sequentially.
func (s *State) ApplyTransaction(txBytes []byte, txIndex int, blockHash common.Hash) error {

	var t ethTypes.Transaction
	if err := rlp.Decode(bytes.NewReader(txBytes), &t); err != nil {
		s.logger.WithError(err).Error("Decoding Transaction")
		return err
	}
	s.logger.WithField("hash", t.Hash().Hex()).Debug("Decoded tx")

	return s.was.ApplyTransaction(t, txIndex, blockHash)
}

func (s *State) CreateGenesisAccounts() error {

	if _, err := os.Stat(s.genesisFile); os.IsNotExist(err) {
		return nil
	}

	contents, err := ioutil.ReadFile(s.genesisFile)
	if err != nil {
		return err
	}

	var genesis struct {
		Alloc bcommon.AccountMap
	}

	var poa struct {
		Poa bcommon.PoaMap
	}

	if err := json.Unmarshal(contents, &genesis); err != nil {
		return err
	}

	for addr, account := range genesis.Alloc {
		address := common.HexToAddress(addr)
		if s.Empty(address) {
			s.was.ethState.AddBalance(address, math.MustParseBig256(account.Balance))
			s.was.ethState.SetCode(address, common.Hex2Bytes(account.Code))
			for key, value := range account.Storage {
				s.was.ethState.SetState(address, common.HexToHash(key), common.HexToHash(value))
			}
			s.logger.WithField("address", addr).Debug("Adding account")

			if account.Authorising {
				s.authorisingAccount = addr
				setPOAFROM(addr)
			}
		}
	}

	if err := json.Unmarshal(contents, &poa); err != nil {
		return err
	}

	if string(poa.Poa.Address) != "" {
		addr := poa.Poa.Address
		balance := poa.Poa.Balance
		abi := poa.Poa.Abi
		code := poa.Poa.Code

		s.logger.Debug("POA CreateGenesisAccounts")

		address := common.HexToAddress(addr)
		if s.Empty(address) {
			s.logger.Debug("POA Create Account")

			s.was.ethState.AddBalance(address, math.MustParseBig256(balance))
			s.was.ethState.SetCode(address, common.Hex2Bytes(code))
			//			for key, value := range account.Storage {
			//				s.was.ethState.SetState(address, common.HexToHash(key), common.HexToHash(value))
			//			}
			s.logger.WithField("address", addr).Debug("Adding account")

			s.authorisingAccount = addr
			setPOAFROM(addr)
			s.authorisingAbi = abi
			setPOAABI(abi)
		}

	}

	if _, err = s.Commit(); err != nil {
		return err
	}

	return nil

}

//Empty reports whether the account is non-existant or empty
func (s *State) Empty(addr common.Address) bool {
	res := s.ethState.Empty(addr)
	s.logger.Debugf("%s Empty? %v", addr.Hex(), res)
	return res
}

//GetBalance returns an account's balance from the main ethState
func (s *State) GetBalance(addr common.Address) *big.Int {
	return s.ethState.GetBalance(addr)
}

//GetNonce returns an account's nonce from the main ethState
func (s *State) GetNonce(addr common.Address) uint64 {
	return s.ethState.GetNonce(addr)
}

//GetCode returns an account's bytecode from the main ethState
func (s *State) GetCode(addr common.Address) []byte {
	return s.ethState.GetCode(addr)
}

//GetPoolNonce returns an account's nonce from the txpool's ethState
func (s *State) GetPoolNonce(addr common.Address) uint64 {
	return s.txPool.ethState.GetNonce(addr)
}

//GetTransaction fetches transactions by hash directly from the DB.
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

//GetReceipt fetches transaction receipts by transaction hash directly from the
//DB
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

//GetGasLimit returns the gas limit set between commit calls
func (s *State) GetGasLimit() uint64 {
	return s.gasLimit
}

//GetAuthorisingAccount returns the address of the smart contract which handles
//the list of authorized peers
func (s *State) GetAuthorisingAccount() string {
	return s.authorisingAccount
}

//GetAuthorisingAbi returns the abi of the smart contract which handles
//the list of authorized peers
func (s *State) GetAuthorisingAbi() string {
	return s.authorisingAbi
}

//CheckAuthorised queries the POA smart-contract to check if the address is
//authorised
func (s *State) CheckAuthorised(addr common.Address) (bool, error) {
	contractAddress := common.HexToAddress(s.GetAuthorisingAccount())

	callData, err := POAABI.Pack("checkAuthorised", addr)
	if err != nil {
		s.logger.Warningf("couldn't pack arguments: %v", err)
	}

	// Apply an ethereum call message (no state update) to query the
	// smart-contract. Since there's no nonce check and the gas price is set to
	// zero, an arbitrary address can be used as the source of the tx.
	ethMsg := ethTypes.NewMessage(POAFROM,
		&contractAddress,
		uint64(1),
		big.NewInt(0),
		s.GetGasLimit(),
		big.NewInt(0),
		callData,
		false)

	s.logger.WithFields(logrus.Fields{
		"addr":     addr.Hex(),
		"callData": hex.EncodeToString(callData),
		"contract": contractAddress.Hex(),
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

//------------------------------------------------------------------------------

// getFdLimit retrieves the number of file descriptors allowed to be opened by this
// process.
func getFdLimit() (int, error) {
	var limit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit); err != nil {
		return 0, err
	}
	return int(limit.Cur), nil
}
