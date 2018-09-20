package engine

import (
	"github.com/mosaicnetworks/evm-lite/src/config"
	"github.com/mosaicnetworks/evm-lite/src/consensus"
	"github.com/mosaicnetworks/evm-lite/src/service"
	"github.com/mosaicnetworks/evm-lite/src/state"
	"github.com/sirupsen/logrus"
)

//Engine is the actor that coordinates State, Service and Consensus
type Engine struct {
	state     *state.State
	service   *service.Service
	consensus consensus.Consensus
}

//NewEngine instantiates a new Engine with coupled State, Service, and Consensus
func NewEngine(config config.Config,
	consensus consensus.Consensus,
	logger *logrus.Logger) (*Engine, error) {
	submitCh := make(chan []byte)

	state, err := state.NewState(logger,
		config.Eth.DbFile,
		config.Eth.Cache)
	if err != nil {
		return nil, err
	}

	service := service.NewService(config.Eth.Genesis,
		config.Eth.Keystore,
		config.Eth.EthAPIAddr,
		config.Eth.PwdFile,
		state,
		submitCh,
		logger)

	if err := consensus.Init(state, service); err != nil {
		return nil, err
	}

	engine := &Engine{
		state:     state,
		service:   service,
		consensus: consensus,
	}

	return engine, nil
}

//Run starts the engine's Service asynchronously and starts the Consensus system
//synchronously
func (e *Engine) Run() error {

	go e.service.Run()

	e.consensus.Run()

	return nil
}
