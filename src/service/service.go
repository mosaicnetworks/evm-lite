package service

import (
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
	"github.com/mosaicnetworks/evm-lite/src/state"
	"github.com/sirupsen/logrus"
)

var defaultGas = uint64(90000)

type infoCallback func() (map[string]string, error)

type Service struct {
	sync.Mutex
	state    *state.State
	submitCh chan []byte
	apiAddr  string
	getInfo  infoCallback
	logger   *logrus.Logger
}

func NewService(apiAddr string,
	state *state.State,
	submitCh chan []byte,
	logger *logrus.Logger) *Service {

	return &Service{
		apiAddr:  apiAddr,
		state:    state,
		submitCh: submitCh,
		logger:   logger}
}

func (m *Service) Run() {
	m.logger.Info("serving api...")
	m.serveAPI()
}

func (m *Service) GetSubmitCh() chan []byte {
	return m.submitCh
}

func (m *Service) SetInfoCallback(f infoCallback) {
	m.getInfo = f
}

func (m *Service) serveAPI() {

	serverMuxEVM := http.NewServeMux()

	r := mux.NewRouter()
	r.HandleFunc("/account/{address}", m.makeHandler(accountHandler)).Methods("GET")
	r.HandleFunc("/call", m.makeHandler(callHandler)).Methods("POST")
	r.HandleFunc("/rawtx", m.makeHandler(rawTransactionHandler)).Methods("POST")
	r.HandleFunc("/tx/{tx_hash}", m.makeHandler(transactionReceiptHandler)).Methods("GET")
	r.HandleFunc("/info", m.makeHandler(infoHandler)).Methods("GET")
	r.HandleFunc("/poa", m.makeHandler(poaHandler)).Methods("GET")
	r.HandleFunc("/genesis", m.makeHandler(genesisHandler)).Methods("GET")

	serverMuxEVM.Handle("/", r)

	m.logger.WithField("apiAddr", m.apiAddr).Debug("EVM-Lite Service serving")
	http.ListenAndServe(m.apiAddr, serverMuxEVM)
}

func (m *Service) makeHandler(fn func(http.ResponseWriter, *http.Request, *Service)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.Lock()
		fn(w, r, m)
		m.Unlock()
	}
}

func (m *Service) checkErr(err error) {
	if err != nil {
		m.logger.WithError(err).Error("ERROR")
		os.Exit(1)
	}
}
