package service

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/gorilla/mux"
	"github.com/mosaicnetworks/evm-lite/src/state"
	"github.com/sirupsen/logrus"
)

var defaultGas = uint64(90000)

type infoCallback func() (map[string]string, error)

type Service struct {
	sync.Mutex
	state       *state.State
	submitCh    chan []byte
	keystoreDir string
	apiAddr     string
	keyStore    *keystore.KeyStore
	pwdFile     string
	getInfo     infoCallback
	logger      *logrus.Logger
}

func NewService(keystoreDir, apiAddr, pwdFile string,
	state *state.State,
	submitCh chan []byte,
	logger *logrus.Logger) *Service {
	return &Service{
		keystoreDir: keystoreDir,
		apiAddr:     apiAddr,
		pwdFile:     pwdFile,
		state:       state,
		submitCh:    submitCh,
		logger:      logger}
}

func (m *Service) Run() {
	m.checkErr(m.makeKeyStore())

	m.checkErr(m.unlockAccounts())

	m.logger.Info("serving api...")
	m.serveAPI()
}

func (m *Service) GetSubmitCh() chan []byte {
	return m.submitCh
}

func (m *Service) SetInfoCallback(f infoCallback) {
	m.getInfo = f
}

func (m *Service) makeKeyStore() error {

	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP

	if err := os.MkdirAll(m.keystoreDir, 0700); err != nil {
		return err
	}

	m.keyStore = keystore.NewKeyStore(m.keystoreDir, scryptN, scryptP)

	return nil
}

func (m *Service) unlockAccounts() error {

	if len(m.keyStore.Accounts()) == 0 {
		return nil
	}

	pwd, err := m.readPwd()
	if err != nil {
		m.logger.WithError(err).Error("Reading PwdFile")
		return err
	}

	for _, ac := range m.keyStore.Accounts() {
		if err := m.keyStore.Unlock(ac, string(pwd)); err != nil {
			return err
		}
		m.logger.WithField("address", ac.Address.Hex()).Debug("Unlocked account")
	}
	return nil
}

func (m *Service) serveAPI() {
	r := mux.NewRouter()
	r.HandleFunc("/account/{address}", m.makeHandler(accountHandler)).Methods("GET")
	r.HandleFunc("/accounts", m.makeHandler(accountsHandler)).Methods("GET")
	r.HandleFunc("/call", m.makeHandler(callHandler)).Methods("POST")
	r.HandleFunc("/tx", m.makeHandler(transactionHandler)).Methods("POST")
	r.HandleFunc("/rawtx", m.makeHandler(rawTransactionHandler)).Methods("POST")
	r.HandleFunc("/tx/{tx_hash}", m.makeHandler(transactionReceiptHandler)).Methods("GET")
	r.HandleFunc("/info", m.makeHandler(infoHandler)).Methods("GET")
	r.HandleFunc("/html/info", m.makeHandler(htmlInfoHandler)).Methods("GET")
	r.HandleFunc("/contract", m.makeHandler(contractHandler)).Methods("GET")
	r.HandleFunc("/poa", m.makeHandler(poaHandler)).Methods("GET")

	http.Handle("/", &CORSServer{r})
	http.ListenAndServe(m.apiAddr, nil)
}

type CORSServer struct {
	r *mux.Router
}

func (s *CORSServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
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

func (m *Service) readPwd() (pwd string, err error) {
	text, err := ioutil.ReadFile(m.pwdFile)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(text), "\n")
	// Sanitise DOS line endings.
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], "\r")
	}
	return lines[0], nil
}
