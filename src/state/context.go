package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

var (
	//CustomChainConfig accounts for the fact that EVM-Lite doesn't really have
	//a concetp of Chain, or Blocks (being consensus-agnostic). The EVM is
	//tightly coupled with this ChainConfig object (cf interpreter.go), so this
	//is a workaround that treats all blocks the same.
	CustomChainConfig = params.ChainConfig{
		ChainID:             big.NewInt(1),
		ConstantinopleBlock: big.NewInt(0),
	}
)

func NewContext(origin common.Address,
	coinbase common.Address,
	gasLimit uint64,
	gasPrice *big.Int) vm.Context {

	context := vm.Context{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     func(uint64) common.Hash { return common.Hash{} },
		// Message information
		Origin:   origin,
		Coinbase: coinbase,
		GasLimit: gasLimit,
		GasPrice: gasPrice,
		//The vm has a dependency on this
		//Anything greate than ConstantinopleBlock will do here.
		BlockNumber: CustomChainConfig.ConstantinopleBlock,
	}

	return context
}
