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
