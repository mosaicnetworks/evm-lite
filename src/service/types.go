package service

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

//JSONAccount
type JSONAccount struct {
	Address string            `json:"address"`
	Balance *big.Int          `json:"balance"`
	Nonce   uint64            `json:"nonce"`
	Storage map[string]string `json:"storage,omitempty"`
	Code    string            `json:"bytecode,omitempty"`
}

type JSONAccountList struct {
	Accounts []JSONAccount `json:"accounts"`
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

type JsonContract struct {
	Address common.Address `json:"address"`
	ABI     string         `json:"abi"`
}

type JsonContractList struct {
	Contracts []JsonContract `json:"contracts"`
}
