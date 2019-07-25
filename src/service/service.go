package service

import (
	"net/http"
	"os"
	"sync"

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
	m.logger.WithField("bind_address", m.apiAddr).Debug("Starting EVM-Lite API service")
	m.serveAPI()
}

func (m *Service) GetSubmitCh() chan []byte {
	return m.submitCh
}

func (m *Service) SetInfoCallback(f infoCallback) {
	m.getInfo = f
}

// Serve registers the API handlers with the DefaultServerMux of the http
// package. It calls ListenAndServe but does not process errors returned by it.
// This is because we do not want to throw an error when the consensus system is
// used in-mem and wants to expose its API on the same endpoint (address:port)
// EVM-Lite.
func (m *Service) serveAPI() {
	// Add handlers to DefaultServerMux
	http.HandleFunc("/account/", m.makeHandler(accountHandler))
	http.HandleFunc("/call", m.makeHandler(callHandler))
	http.HandleFunc("/rawtx", m.makeHandler(rawTransactionHandler))
	http.HandleFunc("/tx/", m.makeHandler(transactionReceiptHandler))
	http.HandleFunc("/info", m.makeHandler(infoHandler))
	http.HandleFunc("/poa", m.makeHandler(poaHandler))
	http.HandleFunc("/genesis", m.makeHandler(genesisHandler))

	// It is possible that another server, started in the same process, is
	// simultaneously using the DefaultServerMux. In which case, the handlers
	// will be accessible from both servers.
	http.ListenAndServe(m.apiAddr, nil)
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
