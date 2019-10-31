package state

import (
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	ethState "github.com/ethereum/go-ethereum/core/state"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/mosaicnetworks/evm-lite/src/currency"
)

var _receiptsPrefix = []byte("receipts-")

// BaseState is a THREAD-SAFE wrapper around a StateDB. It contains the logic
// to retrieve information from the DB, and apply new transactions.
type BaseState struct {
	sync.Mutex
	db           ethdb.Database
	stateDB      *ethState.StateDB
	signer       ethTypes.Signer
	chainConfig  params.ChainConfig
	vmConfig     vm.Config
	gasLimit     uint64
	gp           *core.GasPool
	totalUsedGas uint64
}

// NewBaseState returns a BaseState initialized from a database and root hash.
func NewBaseState(db ethdb.Database,
	root common.Hash,
	signer ethTypes.Signer,
	chainConfig params.ChainConfig,
	vmConfig vm.Config,
	gasLimit uint64) BaseState {

	stateDB, _ := ethState.New(root, ethState.NewDatabase(db))

	return BaseState{
		db:          db,
		stateDB:     stateDB,
		signer:      signer,
		chainConfig: chainConfig,
		vmConfig:    vmConfig,
		gasLimit:    gasLimit,
		gp:          new(core.GasPool).AddGas(gasLimit),
	}
}

// Copy returns a copy of the BaseState with its own mutex, gp, and stateDB
func (bs *BaseState) Copy() BaseState {
	return BaseState{
		db:          bs.db,
		stateDB:     bs.stateDB.Copy(),
		signer:      bs.signer,
		chainConfig: bs.chainConfig,
		vmConfig:    bs.vmConfig,
		gasLimit:    bs.gasLimit,
		gp:          new(core.GasPool).AddGas(bs.gasLimit),
	}
}

// CreateAccount adds an account in the stateDB
func (bs *BaseState) CreateAccount(address common.Address,
	code string,
	storage map[string]string,
	balance string,
	nonce uint64) {

	bs.Lock()
	defer bs.Unlock()

	if bs.stateDB.Empty(address) {
		bs.stateDB.AddBalance(address, math.MustParseBig256(currency.ExpandCurrencyString(balance)))
		bs.stateDB.SetCode(address, common.Hex2Bytes(code))
		for key, value := range storage {
			bs.stateDB.SetState(address, common.HexToHash(key), common.HexToHash(value))
		}
		bs.stateDB.SetNonce(address, nonce)
	}
}

// ApplyTransaction executes the transaction on the stateDB and sets its receipt
// unless noReceipt is set to true. An error indicates a consensus issue.
func (bs *BaseState) ApplyTransaction(
	tx *EVMLTransaction,
	txIndex int,
	blockHash common.Hash,
	coinbase common.Address,
	noReceipt bool) error {

	msg := tx.Msg()

	context := NewContext(msg.From(), coinbase, msg.Gas(), msg.GasPrice())

	bs.Lock()
	defer bs.Unlock()

	// Prepare the stateDB with transaction Hash so that it can be used in
	// emitted logs. Not required for CheckTx with no receipt produced.
	if !noReceipt {
		bs.stateDB.Prepare(tx.Hash(), blockHash, txIndex)
	}

	vmenv := vm.NewEVM(context, bs.stateDB, &bs.chainConfig, bs.vmConfig)

	// Apply the transaction to the stateDB (included in the env)
	_, gas, failed, err := core.ApplyMessage(vmenv, msg, bs.gp)
	if (err != nil) || (noReceipt) {
		// These are called "consensus" errors. Return immediately.
		return err
	}

	bs.totalUsedGas += gas

	// Compute the current root hash of the state trie, which will go in the
	// receipt. This has side effects; it updates StateObjects like
	// smart-contract memory.
	root := bs.stateDB.IntermediateRoot(true)

	receipt := ethTypes.NewReceipt(root.Bytes(), failed, bs.totalUsedGas)
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = gas

	// if the transaction created a contract, store the creation address in the
	// receipt.
	if msg.To() == nil {
		receipt.ContractAddress = crypto.CreateAddress(vmenv.Context.Origin, tx.Nonce())
	}

	// Set the receipt logs
	receipt.Logs = bs.stateDB.GetLogs(tx.Hash())

	// set the EVMLTransaction's receipt
	tx.receipt = receipt

	return nil
}

// Call executes a readonly transaction on a copy of the stateDB. This is used
// to call smart-contract methods. We use a copy of the stateDB because even a
// call transaction increments the sender's nonce.
func (bs *BaseState) Call(callMsg ethTypes.Message) ([]byte, error) {
	bs.Lock()
	defer bs.Unlock()

	context := NewContext(callMsg.From(), common.Address{}, 0, big.NewInt(0))

	vmenv := vm.NewEVM(context, bs.stateDB.Copy(), &bs.chainConfig, bs.vmConfig)

	res, _, _, err := core.ApplyMessage(vmenv, callMsg, new(core.GasPool).AddGas(bs.gasLimit))

	return res, err
}

// Reset resets the stateDB and the gas counters.
func (bs *BaseState) Reset(root common.Hash) error {
	bs.Lock()
	defer bs.Unlock()

	err := bs.stateDB.Reset(root)
	if err != nil {
		return err
	}

	bs.totalUsedGas = 0
	bs.gp = new(core.GasPool).AddGas(bs.gasLimit)

	return nil
}

// Commit commits everything to the underlying database
func (bs *BaseState) Commit() (common.Hash, error) {
	bs.Lock()
	bs.Unlock()

	root, err := bs.stateDB.Commit(true)
	if err != nil {
		return common.Hash{}, err
	}

	// FORCE DISK WRITE
	// Apparenty Geth does something smarter here, but can't figure it out
	bs.stateDB.Database().TrieDB().Commit(root, true)

	return root, nil
}

// GetBalance returns an account's balance from the stateDB
func (bs *BaseState) GetBalance(addr common.Address) *big.Int {
	bs.Lock()
	defer bs.Unlock()
	return bs.stateDB.GetBalance(addr)
}

// GetNonce returns an account's nonce from the stateDB
func (bs *BaseState) GetNonce(addr common.Address) uint64 {
	bs.Lock()
	defer bs.Unlock()
	return bs.stateDB.GetNonce(addr)
}

// GetCode returns an account's bytecode from the stateDB
func (bs *BaseState) GetCode(addr common.Address) []byte {
	bs.Lock()
	defer bs.Unlock()
	return bs.stateDB.GetCode(addr)
}

// GetStorage returns an account's storage  from the stateDB
func (bs *BaseState) GetStorage(addr common.Address) map[string]string {
	bs.Lock()
	defer bs.Unlock()

	//	func (db *StateDB) ForEachStorage(addr common.Address, cb func(key, value common.Hash) bool) {
	storage := make(map[string]string)

	bs.stateDB.ForEachStorage(addr, func(key, value common.Hash) bool {
		storage[strings.TrimPrefix(key.Hex(), "0x")] = strings.TrimPrefix(value.Hex(), "0x")
		return true
	})

	return storage
}

// WriteTransactions writes a set of transactions directly into the DB
func (bs *BaseState) WriteTransactions(txs map[common.Hash]*EVMLTransaction) error {
	batch := bs.db.NewBatch()

	for hash, tx := range txs {
		if err := batch.Put(hash.Bytes(), tx.rlpBytes); err != nil {
			return err
		}
	}

	// Write the scheduled data into the database
	return batch.Write()
}

// WriteReceipts writes a set of receipts directly into the DB
func (bs *BaseState) WriteReceipts(txs map[common.Hash]*EVMLTransaction) error {
	batch := bs.db.NewBatch()

	for txHash, tx := range txs {
		storageReceipt := (*ethTypes.ReceiptForStorage)(tx.receipt)
		data, err := rlp.EncodeToBytes(storageReceipt)
		if err != nil {
			return err
		}
		if err := batch.Put(append(_receiptsPrefix, txHash.Bytes()...), data); err != nil {
			return err
		}
	}

	return batch.Write()
}

// GetTransaction fetches transactions by hash directly from the DB.
func (bs *BaseState) GetTransaction(hash common.Hash) (*ethTypes.Transaction, error) {
	// Retrieve the transaction itself from the database
	data, err := bs.db.Get(hash.Bytes())
	if err != nil {
		return nil, err
	}
	var tx ethTypes.Transaction
	if err := rlp.DecodeBytes(data, &tx); err != nil {
		return nil, err
	}

	return &tx, nil
}

// GetReceipt fetches transaction receipts by transaction hash directly from the
// DB
func (bs *BaseState) GetReceipt(txHash common.Hash) (*ethTypes.Receipt, error) {
	data, err := bs.db.Get(append(_receiptsPrefix, txHash.Bytes()...))
	if err != nil {
		return nil, err
	}

	var receipt ethTypes.ReceiptForStorage
	if err := rlp.DecodeBytes(data, &receipt); err != nil {
		return nil, err
	}

	return (*ethTypes.Receipt)(&receipt), nil
}
