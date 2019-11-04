package service

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

//JSONAccount is the JSON structure used for the account endpoint
type JSONAccount struct {
	Address string            `json:"address"`
	Balance *big.Int          `json:"balance"`
	Nonce   uint64            `json:"nonce"`
	Storage map[string]string `json:"storage,omitempty"`
	Code    string            `json:"bytecode,omitempty"`
}

/* // This is unused. Commented out pending deletion
type JSONAccountList struct {
	Accounts []JSONAccount `json:"accounts"`
}
*/

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

//JSONCallRes is the JSON structure for the return from the call endpoint
type JSONCallRes struct {
	Data string `json:"data"`
}

//JSONTxRes has been replaced by JSONReceipt
type JSONTxRes struct {
	TxHash string `json:"txHash"`
}

//JSONContract is the JSON structure returned by the poa endpoint
type JSONContract struct {
	Address common.Address `json:"address"`
	ABI     string         `json:"abi"`
}

/* //Not used. Commented out, pending deletion.
type JSONContractList struct {
	Contracts []JSONContract `json:"contracts"`
}
*/
