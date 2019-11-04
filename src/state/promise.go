package state

import (
	ethCommon "github.com/ethereum/go-ethereum/common"

	"github.com/mosaicnetworks/evm-lite/src/common"
)

// ReceiptPromiseResponse captures a receipt and a potential error
type ReceiptPromiseResponse struct {
	Receipt *common.JSONReceipt
	Error   error
}

// ReceiptPromise provides a response mechanism for transaction receipts. The
// Hash identifies the transaction to which the ReceiptPromise corresponds, and
// is used as the key in the map kept by the WAS.
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

// Respond resolves a ReceiptPromiseResponse and passes it to the RespCh
func (p *ReceiptPromise) Respond(receipt *common.JSONReceipt, err error) {
	p.RespCh <- ReceiptPromiseResponse{receipt, err}
}
