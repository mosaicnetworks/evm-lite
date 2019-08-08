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
	logger   *logrus.Entry
}

func NewService(apiAddr string,
	state *state.State,
	submitCh chan []byte,
	logger *logrus.Entry) *Service {

	return &Service{
		apiAddr:  apiAddr,
		state:    state,
		submitCh: submitCh,
		logger:   logger,
		// rcptproms: make(map[common.Hash]ReceiptPromise),
	}
}

func (m *Service) Run() {
	m.logger.WithField("bind_address", m.apiAddr).Info("API")
	m.serveAPI()
}

func (m *Service) GetSubmitCh() chan []byte {
	return m.submitCh
}

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

	// The call to ListenAndServe is a blocking operation
	err := http.ListenAndServe(m.apiAddr, nil)
	if err != nil {
		m.logger.Error(err)
	}
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
