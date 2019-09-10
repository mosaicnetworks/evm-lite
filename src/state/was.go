package state

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethState "github.com/ethereum/go-ethereum/core/state"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	bcommon "github.com/mosaicnetworks/evm-lite/src/common"
	"github.com/sirupsen/logrus"
)

// WriteAheadState is a wrapper around a DB and StateDB object that applies
// transactions to the StateDB and only commits them to the DB upon Commit. It
// also handles persisting transactions, logs, and receipts to the DB.
// NOT THREAD SAFE
type WriteAheadState struct {
	db       ethdb.Database
	ethState *ethState.StateDB

	signer      ethTypes.Signer
	chainConfig params.ChainConfig // vm.env is still tightly coupled with chainConfig
	vmConfig    vm.Config
	gasLimit    uint64

	txIndex      int
	transactions []*ethTypes.Transaction
	receipts     []*ethTypes.Receipt
	allLogs      []*ethTypes.Log

	totalUsedGas uint64
	gp           *core.GasPool

	logger *logrus.Entry

	receiptPromises map[common.Hash]*ReceiptPromise
}

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
		logger:          logger,
		receiptPromises: make(map[common.Hash]*ReceiptPromise),
	}, nil
}

func (was *WriteAheadState) Reset(root common.Hash) error {

	err := was.ethState.Reset(root)
	if err != nil {
		return err
	}

	was.txIndex = 0
	was.transactions = []*ethTypes.Transaction{}
	was.receipts = []*ethTypes.Receipt{}
	was.allLogs = []*ethTypes.Log{}

	was.totalUsedGas = 0
	was.gp = new(core.GasPool).AddGas(was.gasLimit)

	return nil
}

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
	was.transactions = append(was.transactions, &tx)
	was.receipts = append(was.receipts, receipt)
	was.allLogs = append(was.allLogs, receipt.Logs...)

	was.logger.WithField("hash", tx.Hash().Hex()).Debug("Applied tx to WAS")

	return nil
}

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

	//XXX FORCE DISK WRITE
	//Apparenty Geth does something smarter here... but cant figure it out
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

	for _, tx := range was.transactions {
		data, err := rlp.EncodeToBytes(tx)
		if err != nil {
			return err
		}
		if err := batch.Put(tx.Hash().Bytes(), data); err != nil {
			return err
		}
	}

	// Write the scheduled data into the database
	return batch.Write()
}

func (was *WriteAheadState) writeReceipts() error {
	batch := was.db.NewBatch()

	for _, receipt := range was.receipts {
		storageReceipt := (*ethTypes.ReceiptForStorage)(receipt)
		data, err := rlp.EncodeToBytes(storageReceipt)
		if err != nil {
			return err
		}
		if err := batch.Put(append(receiptsPrefix, receipt.TxHash.Bytes()...), data); err != nil {
			return err
		}
	}

	return batch.Write()
}

func (was *WriteAheadState) respondReceiptPromises() error {
	for _, tx := range was.transactions {
		if promise, ok := was.receiptPromises[tx.Hash()]; ok {
			receipt, err := was.getJSONReceipt(tx.Hash())
			promise.Respond(receipt, err)
			delete(was.receiptPromises, tx.Hash())
		}
	}
	return nil
}

func (was *WriteAheadState) getReceipt(txHash common.Hash) (*ethTypes.Receipt, error) {
	data, err := was.db.Get(append(receiptsPrefix, txHash.Bytes()...))
	if err != nil {
		return nil, fmt.Errorf("Getting Receipt: %v", err)
	}
	var receipt ethTypes.ReceiptForStorage
	if err := rlp.DecodeBytes(data, &receipt); err != nil {
		return nil, fmt.Errorf("Decoding Receipt: %v", err)
	}

	return (*ethTypes.Receipt)(&receipt), nil
}

func (was *WriteAheadState) getTransaction(hash common.Hash) (*ethTypes.Transaction, error) {
	data, err := was.db.Get(hash.Bytes())
	if err != nil {
		return nil, fmt.Errorf("Getting Transaction: %v", err)
	}
	var tx ethTypes.Transaction
	if err := rlp.DecodeBytes(data, &tx); err != nil {
		return nil, fmt.Errorf("Decoding Transaction: %v", err)
	}

	return &tx, nil
}

func (was *WriteAheadState) getJSONReceipt(hash common.Hash) (*bcommon.JsonReceipt, error) {
	tx, err := was.getTransaction(hash)
	if err != nil {
		return nil, err
	}

	receipt, err := was.getReceipt(hash)
	if err != nil {
		return nil, err
	}

	signer := ethTypes.NewEIP155Signer(big.NewInt(1))
	from, err := ethTypes.Sender(signer, tx)
	if err != nil {
		return nil, err
	}

	jsonReceipt := bcommon.JsonReceipt{
		Root:              common.BytesToHash(receipt.PostState),
		TransactionHash:   hash,
		From:              from,
		To:                tx.To(),
		GasUsed:           receipt.GasUsed,
		CumulativeGasUsed: receipt.CumulativeGasUsed,
		ContractAddress:   receipt.ContractAddress,
		Logs:              receipt.Logs,
		LogsBloom:         receipt.Bloom,
		Status:            receipt.Status,
	}

	if receipt.Logs == nil {
		jsonReceipt.Logs = []*ethTypes.Log{}
	}

	return &jsonReceipt, nil
}
