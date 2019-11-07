package service

import (
	"math/big"
	"net/http"
	"os"

	"github.com/mosaicnetworks/evm-lite/src/state"
	"github.com/sirupsen/logrus"
)

type infoCallback func() (map[string]string, error)

//Service controls the EVM-Lite endpoints
type Service struct {
	state       *state.State
	submitCh    chan []byte
	apiAddr     string
	minGasPrice *big.Int
	getInfo     infoCallback
	logger      *logrus.Entry
}

//NewService is a factory method that returns a new instance of Service
func NewService(apiAddr string,
	state *state.State,
	submitCh chan []byte,
	minGasPrice *big.Int,
	logger *logrus.Entry) *Service {

	return &Service{
		apiAddr:     apiAddr,
		state:       state,
		submitCh:    submitCh,
		minGasPrice: minGasPrice,
		logger:      logger,
	}
}

//Run starts the Service serving
func (m *Service) Run() {
	m.logger.WithField("bind_address", m.apiAddr).Info("API")
	m.serveAPI()
}

//GetSubmitCh returns the submit channel
func (m *Service) GetSubmitCh() chan []byte {
	return m.submitCh
}

//SetInfoCallback the info callback. This is used to call the babble stats
//endpoint to populate the info endpoint
func (m *Service) SetInfoCallback(f infoCallback) {
	m.getInfo = f
}

// Serve registers the API handlers with the DefaultServerMux of the http
// package, and calls ListenAndServe. It is possible that another module in the
// application (ex: the consensus system) has registered other handlers with the
// DefaultServeMux. In this case, those handlers will also be process by this
// server.
func (m *Service) serveAPI() {
	// Add handlers to DefaultServerMux
	http.HandleFunc("/account/", m.makeHandler(accountHandler))
	http.HandleFunc("/call", m.makeHandler(callHandler))
	http.HandleFunc("/rawtx", m.makeHandler(rawTransactionHandler))
	http.HandleFunc("/tx/", m.makeHandler(transactionReceiptHandler))
	http.HandleFunc("/info", m.makeHandler(infoHandler))
	http.HandleFunc("/poa", m.makeHandler(poaHandler))
	http.HandleFunc("/genesis", m.makeHandler(genesisHandler))
	http.HandleFunc("/version", m.makeHandler(versionHandler))

	//TODO - this is experimental and placed on an endpoint for convenience.
	http.HandleFunc("/export", m.makeHandler(exportHandler))

	// The call to ListenAndServe is a blocking operation
	err := http.ListenAndServe(m.apiAddr, nil)
	if err != nil {
		m.logger.Error(err)
	}
}

func (m *Service) makeHandler(fn func(http.ResponseWriter, *http.Request, *Service)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		fn(w, r, m)
	}
}

func (m *Service) checkErr(err error) {
	if err != nil {
		m.logger.WithError(err).Error("ERROR")
		os.Exit(1)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
