package state

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	eth_common "github.com/ethereum/go-ethereum/common"
)

/*
The Controller smart-contract needs to implement a very simple interface. The only
required method is: ControllerContractAddress returning address
*/

const defaultControllerABI = "[ {  \"constant\": true,  \"inputs\": [],  \"name\": \"POAContractAddress\",  \"outputs\": [   {    \"internalType\": \"address\",    \"name\": \"contractAddress\",    \"type\": \"address\"   }  ],  \"payable\": false,  \"stateMutability\": \"view\",  \"type\": \"function\" }]"
const defaultControllerADDR = "0XAABBAABBAABBAABBAABBAABBAABBAABBAABBAABB"

var (
	// ControllerABI defines the ABI of the Controller smart-contract as needed by a consensus
	// module to check if an address is authorized
	ControllerABI abi.ABI

	// ControllerABISTRING is the string representaion of ControllerABI
	ControllerABISTRING string

	// ControllerADDR is the address of the Controller smart-contract.
	ControllerADDR eth_common.Address
)

func init() {
	ControllerABI, _ = abi.JSON(strings.NewReader(defaultControllerABI))
	ControllerADDR = eth_common.HexToAddress(defaultControllerADDR)
}

func setControllerABI(poaABI string) {
	ControllerABISTRING = poaABI
	ControllerABI, _ = abi.JSON(strings.NewReader(poaABI))
}

func setControllerADDR(address string) {
	ControllerADDR = eth_common.HexToAddress(address)
}
