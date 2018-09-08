package solo

import (
	"github.com/mosaicnetworks/babble/hashgraph"
	"github.com/mosaicnetworks/evm-lite/service"
	"github.com/mosaicnetworks/evm-lite/state"
	"github.com/sirupsen/logrus"
)

/*
Solo implements the Consensus interface.
It relays messages directly from the State to the Service.
*/
type Solo struct {
	blockIndex int
	state      *state.State
	service    *service.Service
	logger     *logrus.Entry
}

//NewSolo returns a Solo object with nil State and Service
func NewSolo(logger *logrus.Logger) *Solo {
	return &Solo{
		logger: logger.WithField("module", "solo"),
	}
}

/*******************************************************************************
IMPLEMENT CONSENSUS INTERFACE
*******************************************************************************/

//Init sets the state and service
func (s *Solo) Init(state *state.State, service *service.Service) error {

	s.logger.Debug("INIT")

	s.state = state
	s.service = service

	return nil
}

//Run pipes the Services's submitCh to the States's ProcessBlock function. It
//wraps individual transactions into Babble Blocks
func (s *Solo) Run() {
	submitCh := s.service.GetSubmitCh()
	for {
		select {
		case t := <-submitCh:
			s.logger.WithField("block", s.blockIndex).Debug("Adding Transaction")

			block := hashgraph.NewBlock(s.blockIndex, 0, []byte{}, [][]byte{t})

			s.logger.WithField("block", s.blockIndex).Debug("Processing Block")

			hash, err := s.state.ProcessBlock(block)
			if err != nil {
				s.logger.WithField("block", s.blockIndex).WithError(err).Error()
			}

			s.logger.WithField("block", s.blockIndex).Debugf("Result State Hash: %s", hash)

			s.blockIndex++
		}
	}
}
