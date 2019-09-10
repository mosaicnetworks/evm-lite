package state

import (
	ethCommon "github.com/ethereum/go-ethereum/common"

	"github.com/mosaicnetworks/evm-lite/src/common"
)

// ReceiptPromiseResponse capture both a response and a potential error
type ReceiptPromiseResponse struct {
	Receipt *common.JsonReceipt
	Error   error
}

// ReceiptPromise provides a response mechanism for transaction receipts. The
// Hash identifies the transaction to which the ReceiptPromise corresponds, and
// is used as the key in the map kept by the state object.
type ReceiptPromise struct {
	Hash   ethCommon.Hash
	RespCh chan ReceiptPromiseResponse
}

// NewReceiptPromise is a factory method for a ReceiptPromise
func NewReceiptPromise(hash ethCommon.Hash) *ReceiptPromise {
	return &ReceiptPromise{
		Hash:   hash,
		RespCh: make(chan ReceiptPromiseResponse, 1),
	}
}

// Respond handles resolving a JsonReceipt
func (p *ReceiptPromise) Respond(receipt *common.JsonReceipt, err error) {
	p.RespCh <- ReceiptPromiseResponse{receipt, err}
}
