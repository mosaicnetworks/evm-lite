package babble

import (
	_babble "github.com/mosaicnetworks/babble/src/babble"
	"github.com/mosaicnetworks/evm-lite/src/config"
	"github.com/mosaicnetworks/evm-lite/src/service"
	"github.com/mosaicnetworks/evm-lite/src/state"
	"github.com/sirupsen/logrus"
)

// InmemBabble implementes the Consensus interface.
// It uses an inmemory Babble node.
type InmemBabble struct {
	config     *config.BabbleConfig
	babble     *_babble.Babble
	ethService *service.Service
	ethState   *state.State
	logger     *logrus.Logger
}

// NewInmemBabble instantiates a new InmemBabble consensus system
func NewInmemBabble(config *config.BabbleConfig, logger *logrus.Logger) *InmemBabble {
	return &InmemBabble{
		config: config,
		logger: logger,
	}
}

/*******************************************************************************
IMPLEMENT CONSENSUS INTERFACE
*******************************************************************************/

// Init instantiates a Babble inmemory node
func (b *InmemBabble) Init(state *state.State, service *service.Service) error {

	b.logger.Debug("INIT")

	b.ethState = state
	b.ethService = service

	realConfig := b.config.ToRealBabbleConfig(b.logger)
	realConfig.Proxy = NewInmemProxy(state, service, service.GetSubmitCh(), b.logger)

	babble := _babble.NewBabble(realConfig)
	err := babble.Init()
	if err != nil {
		return err
	}
	b.babble = babble

	return nil
}

// Run starts the Babble node
func (b *InmemBabble) Run() error {
	b.babble.Run()
	return nil
}

// Info returns Babble stats
func (b *InmemBabble) Info() (map[string]string, error) {
	info := b.babble.Node.GetStats()
	info["type"] = "babble"
	return info, nil
}
