package state

import (
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

// TxPool is a BaseState extension with a CheckTx function. The service requires
// a stateful object to check transactions, because txs might be coming it
// faster than the consensus system can process them.
type TxPool struct {
	BaseState
	logger *logrus.Entry
}

// NewTxPool creates a new TxPool object
func NewTxPool(base BaseState, logger *logrus.Entry) *TxPool {

	return &TxPool{
		BaseState: base,
		logger:    logger,
	}
}

// CheckTx applies the transaction to the base's stateDB. It doesn't care about
// the transaction index, block hash, or coinbase. It is used by the service to
// quickly check if a transaction is valid before submitting it to the consensus
// system.
func (p *TxPool) CheckTx(tx *ethTypes.Transaction) error {
	_, err := p.ApplyTransaction(*tx, 0, common.Hash{}, common.Address{}, true)
	return err
}
