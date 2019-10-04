package common

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

type Genesis struct {
	Alloc AccountMap
	Poa   PoaMap
}

type AccountMap map[string]struct {
	Code        string
	Storage     map[string]string
	Balance     string
	Authorising bool
}

type PoaMap struct {
	Address string
	Balance string
	Abi     string
	Code    string
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

func ToJSONReceiptNoFrom(receipt *ethTypes.Receipt, tx *ethTypes.Transaction, signer ethTypes.Signer) *JsonReceipt {
	from, _ := ethTypes.Sender(signer, tx)
	return ToJSONReceipt(receipt, tx, signer, from)
}

func ToJSONReceipt(receipt *ethTypes.Receipt, tx *ethTypes.Transaction, signer ethTypes.Signer, from ethcommon.Address) *JsonReceipt {

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
