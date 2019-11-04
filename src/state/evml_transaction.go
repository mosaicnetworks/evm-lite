package state

import (
	"bytes"

	ethCommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/mosaicnetworks/evm-lite/src/common"
)

// EVMLTransaction is a wrapper around an EVM transaction which contains a
// receipt and the sender address.
type EVMLTransaction struct {
	*ethTypes.Transaction
	message  *ethTypes.Message
	receipt  *ethTypes.Receipt
	rlpBytes []byte
}

// NewEVMLTransaction decodes an RLP encoded EVM transaction and returns an
// EVMLTransaction wrapper
func NewEVMLTransaction(rlpBytes []byte, signer ethTypes.Signer) (*EVMLTransaction, error) {

	var t ethTypes.Transaction
	if err := rlp.Decode(bytes.NewReader(rlpBytes), &t); err != nil {
		return nil, err
	}

	msg, err := t.AsMessage(signer)
	if err != nil {
		return nil, err
	}

	evmlt := &EVMLTransaction{
		Transaction: &t,
		message:     &msg,
		rlpBytes:    rlpBytes,
	}

	return evmlt, nil
}

// Msg returns the transaction's core.Message property
func (t *EVMLTransaction) Msg() *ethTypes.Message {
	return t.message
}

// From returns the transaction's sender address
func (t *EVMLTransaction) From() ethCommon.Address {
	if t.message != nil {
		return t.message.From()
	}
	return ethCommon.Address{}
}

// JSONReceipt returns the JSONReceipt corresponding to an EVMLTransaction
func (t *EVMLTransaction) JSONReceipt() *common.JSONReceipt {
	if t.receipt == nil {
		return nil
	}

	jsonReceipt := common.JSONReceipt{
		Root:              ethCommon.BytesToHash(t.receipt.PostState),
		TransactionHash:   t.Hash(),
		From:              t.From(),
		To:                t.To(),
		GasUsed:           t.receipt.GasUsed,
		CumulativeGasUsed: t.receipt.CumulativeGasUsed,
		ContractAddress:   t.receipt.ContractAddress,
		Logs:              t.receipt.Logs,
		LogsBloom:         t.receipt.Bloom,
		Status:            t.receipt.Status,
	}

	if t.receipt.Logs == nil {
		jsonReceipt.Logs = []*ethTypes.Log{}
	}

	return &jsonReceipt
}
