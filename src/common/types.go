package common

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

//Genesis File Structure
type Genesis struct {
	Alloc AccountMap
	Poa   PoaMap
}

//AccountMap holds the alloc section of the genesis file
type AccountMap map[string]struct {
	Code        string            `json:"code"`
	Storage     map[string]string `json:"storage"`
	Balance     string            `json:"balance"`
	Authorising bool              `json:"authorising"`
	Nonce       uint64            `json:"nonce,omitempty"`
}

//PoaMap holds the poa section of the genesis file
type PoaMap struct {
	Address string
	Balance string
	Storage map[string]string
	Abi     string
	Code    string
	Nonce   uint64 `json:"nonce,omitempty"`
}

type JsonReceipt struct {
	Root              ethcommon.Hash     `json:"root"`
	TransactionHash   ethcommon.Hash     `json:"transactionHash"`
	From              ethcommon.Address  `json:"from"`
	To                *ethcommon.Address `json:"to"`
	GasUsed           uint64             `json:"gasUsed"`
	CumulativeGasUsed uint64             `json:"cumulativeGasUsed"`
	ContractAddress   ethcommon.Address  `json:"contractAddress"`
	Logs              []*ethTypes.Log    `json:"logs"`
	LogsBloom         ethTypes.Bloom     `json:"logsBloom"`
	Status            uint64             `json:"status"`
}

// ToJSONReceipt uses a transaction, its from address, and a receipt to create
// a JSONReceipt. The "from" addressed is derived from the transaction's
// signature.
func ToJSONReceipt(receipt *ethTypes.Receipt, tx *ethTypes.Transaction, signer ethTypes.Signer) *JsonReceipt {
	from, _ := ethTypes.Sender(signer, tx)

	jsonReceipt := JsonReceipt{
		Root:              ethcommon.BytesToHash(receipt.PostState),
		TransactionHash:   tx.Hash(),
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

	return &jsonReceipt
}
