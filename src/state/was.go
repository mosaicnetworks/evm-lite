package state

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	evmlCommon "github.com/mosaicnetworks/evm-lite/src/common"
	"github.com/sirupsen/logrus"
)

type evmlTransactions struct {
	transaction *ethTypes.Transaction
	receipt     *ethTypes.Receipt
	from        common.Address
	txBytes     []byte
}

// WriteAheadState is a wrapper around a DB and StateBase object that applies
// transactions to the StateDB and only commits them to the DB upon Commit. It
// also handles persisting transactions, logs, and receipts to the DB.
// NOT THREAD SAFE
type WriteAheadState struct {
	BaseState

	txIndex int

	tx map[common.Hash]*evmlTransactions //replaces transactions, receipts and receiptPromises

	//	transactions map[common.Hash]*ethTypes.Transaction //Deprecated
	//	receipts     map[common.Hash]*ethTypes.Receipt     //Deprecated
	allLogs []*ethTypes.Log

	receiptPromises map[common.Hash]*ReceiptPromise
	promiseLock     sync.Mutex

	logger *logrus.Entry
}

// NewWriteAheadState returns a new WAS based on a BaseState
func NewWriteAheadState(base BaseState, logger *logrus.Entry) *WriteAheadState {
	return &WriteAheadState{
		BaseState: base,
		//		transactions:    make(map[common.Hash]*ethTypes.Transaction),
		//		receipts:        make(map[common.Hash]*ethTypes.Receipt),
		tx:              make(map[common.Hash]*evmlTransactions),
		receiptPromises: make(map[common.Hash]*ReceiptPromise),
		logger:          logger,
	}
}

// Reset overrides BaseState Reset. It calls reset on the BaseState and clears
// the transactions, receipts, and logs caches.
func (was *WriteAheadState) Reset(root common.Hash) error {

	err := was.BaseState.Reset(root)
	if err != nil {
		return err
	}

	was.txIndex = 0
	was.tx = make(map[common.Hash]*evmlTransactions)

	//	was.transactions = make(map[common.Hash]*ethTypes.Transaction)
	//	was.receipts = make(map[common.Hash]*ethTypes.Receipt)
	was.allLogs = []*ethTypes.Log{}

	return nil
}

// CreateReceiptPromise creates and records a new ReceiptPromise for a
// transaction hash.
func (was *WriteAheadState) CreateReceiptPromise(hash common.Hash) *ReceiptPromise {
	was.promiseLock.Lock()
	defer was.promiseLock.Unlock()

	p := NewReceiptPromise(hash)

	was.receiptPromises[hash] = p

	return p
}

// ApplyTransaction executes the transaction on the WAS BaseState. If the
// transaction returns a "consensus" error (an error that is not due to EVM
// execution), it will not produce a receipt, and will not be saved; if there is
// a promise attached to it, we quickly resolve it with an error. If the
// transaction did not return a "consensus" error, we record it and its receipt,
// even if its status is "failed".
func (was *WriteAheadState) ApplyTransaction(
	tx ethTypes.Transaction,
	txIndex int,
	blockHash common.Hash,
	coinbase common.Address,
	txBytes []byte) error {

	txHash := tx.Hash()

	// Apply the transaction to the current state (included in the env)
	receipt, from, err := was.BaseState.ApplyTransaction(tx, txIndex, blockHash, coinbase, false)
	if err != nil {
		was.logger.WithError(err).Error("Applying transaction to WAS")

		// Respond to the promise immediately if we got a "consensus" error
		if promise, ok := was.receiptPromises[txHash]; ok {
			promise.Respond(nil, err)
			delete(was.receiptPromises, txHash)
		}

		return err
	}

	was.txIndex++

	was.tx[txHash] = &evmlTransactions{transaction: &tx, receipt: receipt, from: from, txBytes: txBytes}

	//	was.transactions[txHash] = &tx
	//	was.receipts[txHash] = receipt

	was.allLogs = append(was.allLogs, receipt.Logs...)

	if was.logger.Level > logrus.InfoLevel {
		was.logger.WithField("hash", txHash.Hex()).Debug("Applied tx to WAS")
	}

	return nil
}

// Commit commits everything to the underlying database.
func (was *WriteAheadState) Commit() (common.Hash, error) {
	was.logger.WithFields(logrus.Fields{
		"txs":      was.txIndex,
		"receipts": len(was.tx),
		"logs":     len(was.allLogs),
	}).Info("Commit")

	// Commit all state changes to the database
	root, err := was.BaseState.Commit()
	if err != nil {
		was.logger.WithError(err).Error("Committing state")
		return common.Hash{}, err
	}

	if err := was.BaseState.WriteTransactions(was.tx); err != nil {
		was.logger.WithError(err).Error("Writing txs")
		return common.Hash{}, err
	}

	if err := was.BaseState.WriteReceipts(was.tx); err != nil {
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

func (was *WriteAheadState) respondReceiptPromises() error {
	was.promiseLock.Lock()
	defer was.promiseLock.Unlock()

	for _, obj := range was.tx {
		tx := obj.transaction
		if promise, ok := was.receiptPromises[tx.Hash()]; ok {

			if obj.receipt == nil {
				promise.Respond(nil, fmt.Errorf("No Transaction Receipt"))
			} else {
				promise.Respond(evmlCommon.ToJSONReceipt(obj.receipt, tx, was.BaseState.signer, obj.from), nil)
			}
			delete(was.receiptPromises, tx.Hash())
		}
	}
	return nil
}
