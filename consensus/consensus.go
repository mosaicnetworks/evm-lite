package consensus

import (
	"github.com/mosaicnetworks/evm-lite/service"
	"github.com/mosaicnetworks/evm-lite/state"
)

//Consensus is the interface that abstracts the consensus system
type Consensus interface {
	Init(*state.State, *service.Service) error
	Run() error
}
