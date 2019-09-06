package engine

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/mosaicnetworks/evm-lite/src/config"
	"github.com/mosaicnetworks/evm-lite/src/consensus"
	"github.com/mosaicnetworks/evm-lite/src/currency"
	"github.com/mosaicnetworks/evm-lite/src/service"
	"github.com/mosaicnetworks/evm-lite/src/state"
)

// Engine is the actor that coordinates State, Service and Consensus
type Engine struct {
	state     *state.State
	service   *service.Service
	consensus consensus.Consensus
}

// NewEngine instantiates a new Engine with coupled State, Service, and Consensus
func NewEngine(config config.Config, consensus consensus.Consensus) (*Engine, error) {

	logger := config.Logger()

	submitCh := make(chan []byte)

	state, err := state.NewState(
		config.DbFile,
		config.Cache,
		config.Genesis,
		logger.WithField("component", "state"))

	if err != nil {
		logger.WithError(err).Error("engine.go:NewEngine() state.NewState")
		return nil, err
	}

	minGasPrice, ok := math.ParseBig256(currency.ExpandCurrencyString(config.MinGasPrice))
	if !ok {
		logger.WithField("min-gas-price", config.MinGasPrice).Debug("Could not parse min-gas-price")
		minGasPrice = big.NewInt(0)
	}

	service := service.NewService(
		config.EthAPIAddr,
		state,
		submitCh,
		minGasPrice,
		logger.WithField("component", "service"))

	if err := consensus.Init(state, service); err != nil {
		logger.WithError(err).Error("engine.go:NewEngine() Consensus.Init")
		return nil, err
	}

	service.SetInfoCallback(consensus.Info)

	engine := &Engine{
		state:     state,
		service:   service,
		consensus: consensus,
	}

	return engine, nil
}

// Run starts the engine's Service asynchronously and starts the Consensus
// system synchronously
func (e *Engine) Run() error {

	go e.service.Run()

	e.consensus.Run()

	return nil
}
