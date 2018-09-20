package babble

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/mosaicnetworks/babble/hashgraph"
	"github.com/mosaicnetworks/evm-lite/src/service"
	"github.com/mosaicnetworks/evm-lite/src/state"
	"github.com/sirupsen/logrus"
)

//InmemProxy implements the Babble AppProxy interface
type InmemProxy struct {
	service  *service.Service
	state    *state.State
	submitCh chan []byte
	logger   *logrus.Entry
}

//NewInmemProxy initializes and return a new InmemProxy
func NewInmemProxy(state *state.State,
	service *service.Service,
	submitCh chan []byte,
	logger *logrus.Entry) *InmemProxy {

	return &InmemProxy{
		service:  service,
		state:    state,
		submitCh: submitCh,
		logger:   logger,
	}
}

/*******************************************************************************
Implement Babble AppProxy Interface
*******************************************************************************/

//SubmitCh is the channel through which the Service sends transactions to the
//node.
func (p *InmemProxy) SubmitCh() chan []byte {
	return p.submitCh
}

//CommitBlock commits Block to the State and expects the resulting state hash
func (p *InmemProxy) CommitBlock(block hashgraph.Block) ([]byte, error) {
	p.logger.Debug("CommitBlock")

	blockHashBytes, err := block.Hash()
	blockHash := common.BytesToHash(blockHashBytes)

	for i, tx := range block.Transactions() {
		if err := p.state.ApplyTransaction(tx, i, blockHash); err != nil {
			return []byte{}, err
		}
	}

	hash, err := p.state.Commit()
	if err != nil {
		return []byte{}, err
	}

	return hash.Bytes(), nil
}

//TODO - Implement these two functions
func (p *InmemProxy) GetSnapshot(blockIndex int) ([]byte, error) {
	return []byte{}, nil
}

func (p *InmemProxy) Restore(snapshot []byte) error {
	return nil
}
