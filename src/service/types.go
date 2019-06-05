package service

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

type JsonAccount struct {
	Address string   `json:"address"`
	Balance *big.Int `json:"balance"`
	Nonce   uint64   `json:"nonce"`
	Code    string   `json:"bytecode"`
}

type JsonAccountList struct {
	Accounts []JsonAccount `json:"accounts"`
}

// SendTxArgs represents the arguments to sumbit a new transaction into the transaction pool.
type SendTxArgs struct {
	From     common.Address  `json:"from"`
	To       *common.Address `json:"to"`
	Gas      uint64          `json:"gas"`
	GasPrice *big.Int        `json:"gasPrice"`
	Value    *big.Int        `json:"value"`
	Data     string          `json:"data"`
	Nonce    *uint64         `json:"nonce"`
}

type JsonCallRes struct {
	Data string `json:"data"`
}

type JsonTxRes struct {
	TxHash string `json:"txHash"`
}

type JsonReceipt struct {
	Root              common.Hash     `json:"root"`
	TransactionHash   common.Hash     `json:"transactionHash"`
	From              common.Address  `json:"from"`
	To                *common.Address `json:"to"`
	GasUsed           uint64          `json:"gasUsed"`
	CumulativeGasUsed uint64          `json:"cumulativeGasUsed"`
	ContractAddress   common.Address  `json:"contractAddress"`
	Logs              []*ethTypes.Log `json:"logs"`
	LogsBloom         ethTypes.Bloom  `json:"logsBloom"`
	Status            uint64          `json:"status"`
}

type JsonContract struct {
	Address common.Address `json:"address"`
	ABI     string         `json:"abi"`
}

type JsonContractList struct {
	Contracts []JsonContract `json:"contracts"`
}
