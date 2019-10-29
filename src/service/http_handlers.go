package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/mosaicnetworks/evm-lite/src/state"
	"github.com/sirupsen/logrus"

	comm "github.com/mosaicnetworks/evm-lite/src/common"
)

/*
GET /account/{address}?frompool={true|false|t|f|T|F|1|0|TRUE|FALSE|True|False}
example: /account/0x50bd8a037442af4cdf631495bcaa5443de19685d
returns: JSON JsonAccount

This endpoint returns information about any account, taken by default from the
main state, or on the TxPool's ethState if `frompool=true`.
*/
func accountHandler(w http.ResponseWriter, r *http.Request, m *Service) {
	param := r.URL.Path[len("/account/"):]
	address := common.HexToAddress(param)

	if m.logger.Level > logrus.InfoLevel {
		m.logger.WithField("param", param).Debug("GET account")
		m.logger.WithField("address", address.Hex()).Debug("GET account")
	}

	var fromPool bool

	// check query param `state`
	qs := r.URL.Query().Get("frompool")
	if qs != "" {
		fp, err := strconv.ParseBool(qs)
		if err != nil {
			m.logger.WithError(err).Error("Error converting string to bool")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fromPool = fp
	}

	nonce := m.state.GetNonce(address, fromPool)
	balance := m.state.GetBalance(address, fromPool)
	code := hexutil.Encode(m.state.GetCode(address, fromPool))

	if code == "0x" {
		code = ""
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
	if m.logger.Level > logrus.InfoLevel {
		m.logger.WithField("request", r).Debug("POST call")
	}

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

This is a SYNCHRONOUS request. We wait for the transaction to go through
consensus, and return the corresponding receipt directly.
*/
func rawTransactionHandler(w http.ResponseWriter, r *http.Request, m *Service) {
	if m.logger.Level > logrus.InfoLevel {
		m.logger.WithField("request", r).Debug("POST rawtx")
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		m.logger.WithError(err).Error("Reading request body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sBody := string(body)

	if m.logger.Level > logrus.InfoLevel {
		m.logger.WithField("body", body)
		m.logger.WithField("body (string)", sBody).Debug()
	}

	rawTxBytes, err := hexutil.Decode(sBody)
	if err != nil {
		m.logger.WithError(err).Error("Reading raw tx from request body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if m.logger.Level > logrus.InfoLevel {
		m.logger.WithField("raw tx bytes", rawTxBytes).Debug()
	}

	tx, err := state.NewEVMLTransaction(rawTxBytes, m.state.GetSigner())
	if err != nil {
		m.logger.WithError(err).Error("Decoding Transaction")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if m.logger.Level > logrus.InfoLevel {
		m.logger.WithFields(logrus.Fields{
			"hash":     tx.Hash().Hex(),
			"from":     tx.From(),
			"to":       tx.To(),
			"payload":  fmt.Sprintf("%x", tx.Data()),
			"gas":      tx.Gas(),
			"gasPrice": tx.GasPrice(),
			"nonce":    tx.Nonce(),
			"value":    tx.Value(),
		}).Debug("Service decoded tx")
	}

	// Check if gasPrice is above set limit
	if m.minGasPrice != nil && tx.GasPrice().Cmp(m.minGasPrice) < 0 {
		err := fmt.Errorf("Gasprice too low. Got %v, MIN: %v", tx.GasPrice(), m.minGasPrice)
		m.logger.Debug(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := m.state.CheckTx(tx); err != nil {
		m.logger.WithError(err).Error("Checking Transaction")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	promise := m.state.CreateReceiptPromise(tx.Hash())

	m.logger.Debug("submitting tx")
	m.submitCh <- rawTxBytes
	m.logger.Debug("submitted tx")

	timeout := time.After(15 * time.Second)
	var receipt *comm.JsonReceipt
	var respErr error

	select {
	case resp := <-promise.RespCh:
		if resp.Error != nil {
			respErr = resp.Error
			break
		}
		receipt = resp.Receipt
	case <-timeout:
		respErr = fmt.Errorf("Timeout waiting for transaction to go through consensus")
		break
	}

	if respErr != nil {
		m.logger.Errorf("RespErr:  %v", respErr)
		http.Error(w, respErr.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(receipt)
	if err != nil {
		m.logger.WithError(err).Error("Marshalling JSON Response")
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

	jsonReceipt := comm.ToJSONReceipt(receipt, tx, m.state.GetSigner())

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

	// Add min_gas_price
	if m.minGasPrice != nil {
		stats["min_gas_price"] = m.minGasPrice.String()
	} else {
		stats["min_gas_price"] = "0"
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
		ABI:     m.state.GetAuthorisingABI(),
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

/*
GET /export
returns: JSON Export of current state

This endpoint returns the content of the genesis.json file.
*/
func exportHandler(w http.ResponseWriter, r *http.Request, m *Service) {
	m.logger.Debug("GET export")

	// var genesis []comm.AccountMap

	js := m.state.DumpAllAccounts()

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
