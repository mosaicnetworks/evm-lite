package solo

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mosaicnetworks/evm-lite/src/service"
	"github.com/mosaicnetworks/evm-lite/src/state"
	"github.com/sirupsen/logrus"
)

/*
Solo implements the Consensus interface.
It relays messages directly from the State to the Service.
*/
type Solo struct {
	txIndex int
	state   *state.State
	service *service.Service
	logger  *logrus.Entry
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

//Run pipes the Service's submitCh to the States's ProcessBlock function. It
//wraps individual transactions into Babble Blocks
func (s *Solo) Run() error {
	submitCh := s.service.GetSubmitCh()
	for {
		select {
		case t := <-submitCh:
			s.logger.WithField("tx", s.txIndex).Debug("Adding Transaction")

			err := s.state.ApplyTransaction(t,
				s.txIndex,
				common.StringToHash(fmt.Sprintf("block %d", s.txIndex)))
			if err != nil {
				s.logger.WithField("tx", s.txIndex).WithError(err).Errorf("ApplyTransaction")
			}

			hash, err := s.state.Commit()
			if err != nil {
				s.logger.WithField("tx", s.txIndex).WithError(err).Errorf("Commit")
			}

			s.logger.WithField("tx", s.txIndex).Debugf("Result State Hash: %v", hash)

			s.txIndex++
		}
	}
}
