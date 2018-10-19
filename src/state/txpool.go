package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethState "github.com/ethereum/go-ethereum/core/state"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/sirupsen/logrus"
)

type TxPool struct {
	ethState *ethState.StateDB

	signer       ethTypes.Signer
	chainConfig  params.ChainConfig // vm.env is still tightly coupled with chainConfig
	vmConfig     vm.Config
	gasLimit     uint64
	totalUsedGas uint64
	gp           *core.GasPool

	logger *logrus.Logger
}

func NewTxPool(ethState *ethState.StateDB,
	signer ethTypes.Signer,
	chainConfig params.ChainConfig,
	vmConfig vm.Config,
	gasLimit uint64,
	logger *logrus.Logger) *TxPool {

	return &TxPool{
		ethState:    ethState,
		signer:      signer,
		chainConfig: chainConfig,
		vmConfig:    vmConfig,
		gasLimit:    gasLimit,
		logger:      logger,
	}
}

func (p *TxPool) Reset(root common.Hash) error {

	err := p.ethState.Reset(root)
	if err != nil {
		return err
	}

	p.totalUsedGas = 0
	p.gp = new(core.GasPool).AddGas(p.gasLimit)

	return nil
}

func (p *TxPool) CheckTx(tx *ethTypes.Transaction) error {

	msg, err := tx.AsMessage(p.signer)
	if err != nil {
		p.logger.WithError(err).Error("Converting Transaction to Message")
		return err
	}

	context := vm.Context{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     func(uint64) common.Hash { return common.Hash{} },
		// Message information
		Origin:      msg.From(),
		GasLimit:    msg.Gas(),
		GasPrice:    msg.GasPrice(),
		BlockNumber: big.NewInt(0), // The vm has a dependency on this.
	}

	// The EVM should never be reused and is not thread safe.
	vmenv := vm.NewEVM(context, p.ethState, &p.chainConfig, p.vmConfig)

	// Apply the transaction to the current state (included in the env)
	_, gas, _, err := core.ApplyMessage(vmenv, msg, p.gp)
	if err != nil {
		p.logger.WithError(err).Error("Applying transaction to TxPool")
		return err
	}

	p.totalUsedGas += gas

	return nil
}

func (p *TxPool) GetNonce(addr common.Address) uint64 {
	return p.ethState.GetNonce(addr)
}
