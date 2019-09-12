package state

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethState "github.com/ethereum/go-ethereum/core/state"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	evmlCommon "github.com/mosaicnetworks/evm-lite/src/common"
	"github.com/sirupsen/logrus"
)

// WriteAheadState is a wrapper around a DB and StateDB object that applies
// transactions to the StateDB and only commits them to the DB upon Commit. It
// also handles persisting transactions, logs, and receipts to the DB.
// NOT THREAD SAFE
type WriteAheadState struct {
	db ethdb.Database

	ethState     *ethState.StateDB
	signer       ethTypes.Signer
	chainConfig  params.ChainConfig // vm.env is still tightly coupled with chainConfig
	vmConfig     vm.Config
	gasLimit     uint64
	totalUsedGas uint64
	gp           *core.GasPool

	txIndex      int
	transactions map[common.Hash]*ethTypes.Transaction
	receipts     map[common.Hash]*ethTypes.Receipt
	allLogs      []*ethTypes.Log

	receiptPromises map[common.Hash]*ReceiptPromise

	logger *logrus.Entry
}

// NewWriteAheadState returns a new WAS with its StateDB initialised from db and
// root.
func NewWriteAheadState(db ethdb.Database,
	root common.Hash,
	signer ethTypes.Signer,
	chainConfig params.ChainConfig,
	vmConfig vm.Config,
	gasLimit uint64,
	logger *logrus.Entry) (*WriteAheadState, error) {

	ethState, err := ethState.New(root, ethState.NewDatabase(db))
	if err != nil {
		return nil, err
	}

	return &WriteAheadState{
		db:              db,
		ethState:        ethState,
		signer:          signer,
		chainConfig:     chainConfig,
		vmConfig:        vmConfig,
		gasLimit:        gasLimit,
		gp:              new(core.GasPool).AddGas(gasLimit),
		transactions:    make(map[common.Hash]*ethTypes.Transaction),
		receipts:        make(map[common.Hash]*ethTypes.Receipt),
		receiptPromises: make(map[common.Hash]*ReceiptPromise),
		logger:          logger,
	}, nil
}

// Reset calls reset on the StateDB and clears the transactions, receipts, and
// logs caches. It also resets the gas counters.
func (was *WriteAheadState) Reset(root common.Hash) error {

	err := was.ethState.Reset(root)
	if err != nil {
		return err
	}

	was.txIndex = 0
	was.transactions = make(map[common.Hash]*ethTypes.Transaction)
	was.receipts = make(map[common.Hash]*ethTypes.Receipt)
	was.allLogs = []*ethTypes.Log{}

	was.totalUsedGas = 0
	was.gp = new(core.GasPool).AddGas(was.gasLimit)

	return nil
}

// ApplyTransaction executes the transaction on the WAS ethState. If the
// transaction returns a "consensus" error (an error that is not due to EVM
// execution), it will not produce a receipt, and will not be saved; if there is
// a promise attached to it, we quickly resolve it with an error. If the
// transaction did not return a "consensus" error, we record it and its receipt,
// even if its status is "failed".
func (was *WriteAheadState) ApplyTransaction(
	tx ethTypes.Transaction,
	txIndex int,
	blockHash common.Hash,
	coinbase common.Address) error {

	msg, err := tx.AsMessage(was.signer)
	if err != nil {
		was.logger.WithError(err).Error("Converting Transaction to Message")
		return err
	}

	context := NewContext(msg.From(), coinbase, msg.Gas(), msg.GasPrice())

	//Prepare the ethState with transaction Hash so that it can be used in emitted
	//logs
	was.ethState.Prepare(tx.Hash(), blockHash, txIndex)

	vmenv := vm.NewEVM(context, was.ethState, &was.chainConfig, was.vmConfig)

	// Apply the transaction to the current state (included in the env)
	_, gas, failed, err := core.ApplyMessage(vmenv, msg, was.gp)
	if err != nil {
		was.logger.WithError(err).Error("Applying transaction to WAS")

		// Respond to the promise immediately if we got a "consensus" error
		if promise, ok := was.receiptPromises[tx.Hash()]; ok {
			promise.Respond(nil, err)
			delete(was.receiptPromises, tx.Hash())
		}

		return err
	}

	was.totalUsedGas += gas

	// Create a new receipt for the transaction, storing the intermediate root and gas used by the tx
	// based on the eip phase, we're passing wether the root touch-delete accounts.
	root := was.ethState.IntermediateRoot(true) //this has side effects. It updates StateObjects (SmartContract memory)
	receipt := ethTypes.NewReceipt(root.Bytes(), failed, was.totalUsedGas)
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = gas
	// if the transaction created a contract, store the creation address in the receipt.
	if msg.To() == nil {
		receipt.ContractAddress = crypto.CreateAddress(vmenv.Context.Origin, tx.Nonce())
	}
	// Set the receipt logs and create a bloom for filtering
	receipt.Logs = was.ethState.GetLogs(tx.Hash())
	//receipt.Logs = s.was.state.Logs()
	receipt.Bloom = ethTypes.CreateBloom(ethTypes.Receipts{receipt})

	was.txIndex++
	was.transactions[tx.Hash()] = &tx
	was.receipts[tx.Hash()] = receipt
	was.allLogs = append(was.allLogs, receipt.Logs...)

	was.logger.WithField("hash", tx.Hash().Hex()).Debug("Applied tx to WAS")

	return nil
}

// Commit commits everything to the underlying database.
func (was *WriteAheadState) Commit() (common.Hash, error) {
	was.logger.WithFields(logrus.Fields{
		"txs":      was.txIndex,
		"receipts": len(was.receipts),
		"logs":     len(was.allLogs),
	}).Info("Commit")

	// Commit all state changes to the database
	root, err := was.ethState.Commit(true)
	if err != nil {
		was.logger.WithError(err).Error("Committing state")
		return common.Hash{}, err
	}

	// FORCE DISK WRITE
	// Apparenty Geth does something smarter here, but can't figure it out
	was.ethState.Database().TrieDB().Commit(root, true)

	if err := was.writeTransactions(); err != nil {
		was.logger.WithError(err).Error("Writing txs")
		return common.Hash{}, err
	}

	if err := was.writeReceipts(); err != nil {
		was.logger.WithError(err).Error("Writing receipts")
		return common.Hash{}, err
	}

	// respond to receipts once committed with no errors
	if err := was.respondReceiptPromises(); err != nil {
		was.logger.WithError(err).Error("Responding receipt promises")
		return common.Hash{}, err
	}

	return root, nil
}

func (was *WriteAheadState) writeTransactions() error {
	batch := was.db.NewBatch()

	for hash, tx := range was.transactions {
		data, err := rlp.EncodeToBytes(tx)
		if err != nil {
			return err
		}
		if err := batch.Put(hash.Bytes(), data); err != nil {
			return err
		}
	}

	// Write the scheduled data into the database
	return batch.Write()
}

func (was *WriteAheadState) writeReceipts() error {
	batch := was.db.NewBatch()

	for txHash, receipt := range was.receipts {
		storageReceipt := (*ethTypes.ReceiptForStorage)(receipt)
		data, err := rlp.EncodeToBytes(storageReceipt)
		if err != nil {
			return err
		}
		if err := batch.Put(append(receiptsPrefix, txHash.Bytes()...), data); err != nil {
			return err
		}
	}

	return batch.Write()
}

func (was *WriteAheadState) respondReceiptPromises() error {
	for _, tx := range was.transactions {
		if promise, ok := was.receiptPromises[tx.Hash()]; ok {
			receipt, ok := was.receipts[tx.Hash()]
			if !ok {
				promise.Respond(nil, fmt.Errorf("No Transaction Receipt"))
			} else {
				promise.Respond(evmlCommon.ToJSONReceipt(receipt, tx, was.signer), nil)
			}
			delete(was.receiptPromises, tx.Hash())
		}
	}
	return nil
}
