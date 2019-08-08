package state

import (
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
	"github.com/sirupsen/logrus"

	bcommon "github.com/mosaicnetworks/evm-lite/src/common"
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

	// danu
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
		receiptPromises: make(map[common.Hash]*ReceiptPromise), // Danu
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

func (was *WriteAheadState) ApplyTransaction(tx ethTypes.Transaction, txIndex int, blockHash common.Hash) error {
	msg, err := tx.AsMessage(was.signer)
	if err != nil {
		was.logger.WithError(err).Error("Converting Transaction to Message")
		return err
	}

	context := NewContext(msg.From(), msg.Gas(), msg.GasPrice())

	//Prepare the ethState with transaction Hash so that it can be used in emitted
	//logs
	was.ethState.Prepare(tx.Hash(), blockHash, txIndex)

	vmenv := vm.NewEVM(context, was.ethState, &was.chainConfig, was.vmConfig)

	// Apply the transaction to the current state (included in the env)
	_, gas, failed, err := core.ApplyMessage(vmenv, msg, was.gp)
	if err != nil {
		was.logger.WithError(err).Error("Applying transaction to WAS")
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

	// Danu
	// create json receipt
	// check if hash is in the map
	p, ok := was.receiptPromises[tx.Hash()]
	if ok {
		signer := ethTypes.NewEIP155Signer(big.NewInt(1))
		from, err := ethTypes.Sender(signer, &tx)
		if err != nil {
			was.logger.WithError(err).Error("Getting Tx Sender")
			return err
		}

		jsonReceipt := bcommon.JsonReceipt{
			Root:              common.BytesToHash(receipt.PostState),
			TransactionHash:   receipt.TxHash,
			From:              from,
			To:                tx.To(),
			GasUsed:           receipt.GasUsed,
			CumulativeGasUsed: receipt.CumulativeGasUsed,
			ContractAddress:   receipt.ContractAddress,
			Logs:              receipt.Logs,
			LogsBloom:         receipt.Bloom,
			Status:            receipt.Status,
		}

		p.Respond(jsonReceipt)
	}

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
