package raft

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	_raft "github.com/hashicorp/raft"
	"github.com/mosaicnetworks/evm-lite/src/config"
	"github.com/mosaicnetworks/evm-lite/src/service"
	"github.com/mosaicnetworks/evm-lite/src/state"
	"github.com/sirupsen/logrus"
)

// Raft implements the Consensus interface.
// It uses Hashicorp Raft
type Raft struct {
	config    config.RaftConfig
	service   *service.Service
	fsm       _raft.FSM
	raftNode  *_raft.Raft
	logger    *logrus.Entry
	terminate chan os.Signal
	txIndex   int
}

// NewRaft returns a new Raft object
func NewRaft(config config.RaftConfig, logger *logrus.Logger) *Raft {
	return &Raft{
		config:    config,
		logger:    logger.WithField("module", "raft"),
		terminate: make(chan os.Signal, 1),
	}
}

/*******************************************************************************
IMPLEMENT CONSENSUS INTERFACE
*******************************************************************************/

// Init instantiates a Raft
func (r *Raft) Init(state *state.State, service *service.Service) error {

	r.logger.Debug("INIT")

	r.service = service

	r.fsm = NewFSM(state, r.logger)

	// Initialize raft node

	//TODO: Use r.config
	config := _raft.DefaultConfig()
	config.LocalID = r.config.LocalID

	// Setup Raft communication.
	transport, err := _raft.NewTCPTransport(r.config.NodeAddr,
		nil,
		3,
		10*time.Second,
		os.Stderr)
	if err != nil {
		return err
	}

	// Create the snapshot store. This allows the Raft to truncate the log.
	snapshots, err := _raft.NewFileSnapshotStore(r.config.SnapshotDir, 1, os.Stderr)
	if err != nil {
		return fmt.Errorf("file snapshot store: %s", err)
	}

	// Create the log store and stable store.
	// TODO: Add option for persistent store
	logStore := _raft.NewInmemStore()
	stableStore := _raft.NewInmemStore()

	// Instantiate the Raft systems.
	ra, err := _raft.NewRaft(config, r.fsm, logStore, stableStore, snapshots, transport)
	if err != nil {
		return fmt.Errorf("new raft: %s", err)
	}

	//TODO: We should be using the new dynmamic membership protocol
	configuration, err := _raft.ReadConfigJSON(fmt.Sprintf("%s/peers.json", r.config.RaftDir))
	if err != nil {
		return fmt.Errorf("Unable to create cluster configuration from peers.json: %v", err)
	}
	ra.BootstrapCluster(configuration)

	r.raftNode = ra

	return nil
}

// Run starts the Raft node and service
func (r *Raft) Run() error {

	// Relay submitCh to Raft
	submitCh := r.service.GetSubmitCh()
	signal.Notify(r.terminate, os.Interrupt)
	for {
		select {
		case t := <-submitCh:
			r.logger.WithFields(logrus.Fields{
				"tx":    r.txIndex,
				"state": r.raftNode.State(),
			}).Debug("Adding Transaction")

			if r.raftNode.State() != _raft.Leader {
				r.logger.Debug("NOT LEADER")
				//TODO: Relay message to leader
				break
			}

			f := r.raftNode.Apply(t, r.config.CommitTimeout)
			if err := f.Error(); err != nil {
				r.logger.WithError(err).Error("Applying Raft tx")
				break
			}

			r.txIndex++
		case <-r.terminate:
			r.logger.Debug("Raft exiting")
			return nil
		}
	}
}

// Info returns Raft stats
func (r *Raft) Info() (map[string]string, error) {
	info := r.raftNode.Stats()
	info["type"] = "raft"
	return info, nil
}
