package raft

import (
	"fmt"
	"io"

	_ethCommon "github.com/ethereum/go-ethereum/common"
	_raft "github.com/hashicorp/raft"
	"github.com/mosaicnetworks/evm-lite/state"
	"github.com/sirupsen/logrus"
)

//FSM wraps a state object and implements the Raft FSM interface
type FSM struct {
	state  *state.State
	logger *logrus.Entry
}

//NewFSM returns a new FSM
func NewFSM(state *state.State, logger *logrus.Entry) *FSM {
	return &FSM{
		state:  state,
		logger: logger,
	}
}

/*******************************************************************************
IMPLEMENT RAFT FSM INTERFACE
*******************************************************************************/

//Apply is invoked once a log entry is committed.
//It applies the log data to the state as a transaction.
func (f *FSM) Apply(log *_raft.Log) interface{} {

	f.logger.WithFields(logrus.Fields{
		"index": log.Index,
		"term":  log.Term,
		"type":  log.Type,
		"data":  log.Data,
	}).Debug("Apply")

	if err := f.state.ApplyTransaction(log.Data, int(log.Index), _ethCommon.Hash{}); err != nil {
		f.logger.WithError(err).Error("Error applying transaction")
		return nil
	}

	hash, err := f.state.Commit()
	if err != nil {
		f.logger.WithError(err).Error("Error committing")
		return nil
	}

	return hash.Bytes()
}

//Snapshot is not implemented yet
func (f *FSM) Snapshot() (_raft.FSMSnapshot, error) {
	return nil, fmt.Errorf("Snapshot function not implemented")
}

//Restore is not implemented yet
func (f *FSM) Restore(io.ReadCloser) error {
	return fmt.Errorf("Restore function not implemented")
}
