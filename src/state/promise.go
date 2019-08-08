package state

import (
	ethCommon "github.com/ethereum/go-ethereum/common"

	"github.com/mosaicnetworks/evm-lite/src/common"
)

// ReceiptPromise is a struct for asyncronous response to fetching a receipt
type ReceiptPromise struct {
	Hash ethCommon.Hash

	// response channel
	RespCh *chan common.JsonReceipt
}

// NewReceiptPromise is a factory method for a JoinPromise
func NewReceiptPromise(hash ethCommon.Hash) *ReceiptPromise {
	channel := make(chan common.JsonReceipt)

	return &ReceiptPromise{
		Hash:   hash,
		RespCh: &channel,
	}
}

// Respond handles resolving a JsonReceipt
func (p *ReceiptPromise) Respond(receipt common.JsonReceipt) {
	*p.RespCh <- receipt
}

// ResponseChannel return the response channel for the promise
func (p *ReceiptPromise) ResponseChannel() *chan common.JsonReceipt {
	return p.RespCh
}
