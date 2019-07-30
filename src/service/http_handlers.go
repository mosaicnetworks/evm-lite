package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/sirupsen/logrus"
)

/*
GET /account/{address}
example: /account/0x50bd8a037442af4cdf631495bcaa5443de19685d
returns: JSON JsonAccount

This endpoint should be used to fetch information about any account.
*/
func accountHandler(w http.ResponseWriter, r *http.Request, m *Service) {
	param := r.URL.Path[len("/account/"):]
	m.logger.WithField("param", param).Debug("GET account")
	address := common.HexToAddress(param)
	m.logger.WithField("address", address.Hex()).Debug("GET account")

	balance := m.state.GetBalance(address)
	nonce := m.state.GetNonce(address)
	code := hexutil.Encode(m.state.GetCode(address))
	if code == "0x" {
		code = ""
	} else {
		m.logger.WithField("code", code).Debug("GET account")
	}

	account := JsonAccount{
		Address: address.Hex(),
		Balance: balance,
		Nonce:   nonce,
		Code:    code,
	}

	js, err := json.Marshal(account)
	if err != nil {
		m.logger.WithError(err).Error("Marshaling JSON response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

/*
POST /call
data: JSON SendTxArgs
returns: JSON JsonCallRes

This endpoint allows calling SmartContract code for READONLY operations. These
calls will NOT modify the EVM state.

The data does NOT need to be signed.
*/
func callHandler(w http.ResponseWriter, r *http.Request, m *Service) {
	m.logger.WithField("request", r).Debug("POST call")

	decoder := json.NewDecoder(r.Body)
	var txArgs SendTxArgs
	err := decoder.Decode(&txArgs)
	if err != nil {
		m.logger.WithError(err).Error("Decoding JSON txArgs")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	callMessage, err := prepareCallMessage(txArgs)
	if err != nil {
		m.logger.WithError(err).Error("Converting to CallMessage")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := m.state.Call(*callMessage)
	if err != nil {
		m.logger.WithError(err).Error("Executing Call")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := JsonCallRes{Data: hexutil.Encode(data)}
	js, err := json.Marshal(res)
	if err != nil {
		m.logger.WithError(err).Error("Marshaling JSON response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

/*
POST /rawtx
data: STRING Hex representation of the raw transaction bytes
	  ex: 0xf8620180830f4240946266b0dd0116416b1dacf36...
returns: JSON JsonTxRes

This endpoint allows sending NON-READONLY transactions ALREADY SIGNED. The
client is left to compose a transaction, sign it and RLP encode it. The
resulting bytes, represented as a Hex string is passed to this method to be
forwarded to the EVM.

This is an ASYNCHRONOUS operation and the effect on the State should be verified
by fetching the transaction' receipt.
*/
func rawTransactionHandler(w http.ResponseWriter, r *http.Request, m *Service) {
	m.logger.WithField("request", r).Debug("POST rawtx")

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		m.logger.WithError(err).Error("Reading request body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	m.logger.WithField("body", body)

	sBody := string(body)
	m.logger.WithField("body (string)", sBody).Debug()
	rawTxBytes, err := hexutil.Decode(sBody)
	if err != nil {
		m.logger.WithError(err).Error("Reading raw tx from request body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	m.logger.WithField("raw tx bytes", rawTxBytes).Debug()

	var t ethTypes.Transaction
	if err := rlp.Decode(bytes.NewReader(rawTxBytes), &t); err != nil {
		m.logger.WithError(err).Error("Decoding Transaction")
		return
	}

	m.logger.WithFields(logrus.Fields{
		"hash":     t.Hash().Hex(),
		"to":       t.To(),
		"payload":  fmt.Sprintf("%x", t.Data()),
		"gas":      t.Gas(),
		"gasPrice": t.GasPrice(),
		"nonce":    t.Nonce(),
		"value":    t.Value(),
	}).Debug("Service decoded tx")

	if err := m.state.CheckTx(&t); err != nil {
		m.logger.WithError(err).Error("Checking Transaction")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m.logger.Debug("submitting tx")
	m.submitCh <- rawTxBytes
	m.logger.Debug("submitted tx")

	res := JsonTxRes{TxHash: t.Hash().Hex()}
	js, err := json.Marshal(res)
	if err != nil {
		m.logger.WithError(err).Error("Marshalling JSON response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

/*
GET /tx/{tx_hash}
ex: /tx/0xbfe1aa80eb704d6342c553ac9f423024f448f7c74b3e38559429d4b7c98ffb99
returns: JSON JsonReceipt

This endpoint allows to retrieve the EVM receipt of a specific transactions if it
exists. When a transaction is applied to the EVM , a receipt is saved to allow
checking if/how the transaction affected the state. This is where one can see such
information as the address of a newly created contract, how much gas was use and
the EVM Logs produced by the execution of the transaction.
*/
func transactionReceiptHandler(w http.ResponseWriter, r *http.Request, m *Service) {
	param := r.URL.Path[len("/tx/"):]
	txHash := common.HexToHash(param)
	m.logger.WithField("tx_hash", txHash.Hex()).Debug("GET tx")

	tx, err := m.state.GetTransaction(txHash)
	if err != nil {
		m.logger.WithError(err).Error("Getting Transaction")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	receipt, err := m.state.GetReceipt(txHash)
	if err != nil {
		m.logger.WithError(err).Error("Getting Receipt")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	signer := ethTypes.NewEIP155Signer(big.NewInt(1))
	from, err := ethTypes.Sender(signer, tx)
	if err != nil {
		m.logger.WithError(err).Error("Getting Tx Sender")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonReceipt := JsonReceipt{
		Root:              common.BytesToHash(receipt.PostState),
		TransactionHash:   txHash,
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

	js, err := json.Marshal(jsonReceipt)
	if err != nil {
		m.logger.WithError(err).Error("Marshaling JSON response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

/*
GET /info
returns: JSON (depends on underlying consensus system)

Info returns information about the consensus system. Each consensus system that
plugs into evm-lite must implement an Info function.
*/
func infoHandler(w http.ResponseWriter, r *http.Request, m *Service) {
	m.logger.Debug("GET info")

	stats, err := m.getInfo()
	if err != nil {
		m.logger.WithError(err).Error("Getting Info")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(stats)
	if err != nil {
		m.logger.WithError(err).Error("Marshaling JSON response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

/*
GET /poa
returns: JsonContract
Returns details of the poa smart contract . Replaces /contract
*/
func poaHandler(w http.ResponseWriter, r *http.Request, m *Service) {
	m.logger.Debug("GET poa")

	var al JsonContract

	al = JsonContract{
		Address: common.HexToAddress(m.state.GetAuthorisingAccount()),
		ABI:     m.state.GetAuthorisingAbi(),
	}

	js, err := json.Marshal(al)
	if err != nil {
		m.logger.WithError(err).Error("Marshaling JSON response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

/*
GET /genesis
returns: JSON Genesis

This endpoint returns the content of the genesis.json file.
*/
func genesisHandler(w http.ResponseWriter, r *http.Request, m *Service) {
	m.logger.Debug("GET genesis")

	genesis, err := m.state.GetGenesis()
	if err != nil {
		m.logger.WithError(err).Error("Getting Genesis")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(genesis)
	if err != nil {
		m.logger.WithError(err).Error("Marshaling JSON response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

//------------------------------------------------------------------------------
func prepareCallMessage(args SendTxArgs) (*ethTypes.Message, error) {

	// Create Call Message
	// Set gasPrice and value to 0 because this is a readonly operation
	msg := ethTypes.NewMessage(args.From,
		args.To,
		0,
		big.NewInt(0),
		args.Gas,
		big.NewInt(0),
		common.FromHex(args.Data),
		false)

	return &msg, nil

}
