package state

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	eth_common "github.com/ethereum/go-ethereum/common"
)

/*
The POA smart-contract needs to implement a very simple interface. The only
required method is: checkAuthorised(address) bool
*/

const defaultPOAABI = "[{\"type\":\"function\",\"inputs\": [{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"checkAuthorised\",\"outputs\": [{\"name\":\"\",\"type\":\"bool\"}]}]"
const defaultPOAADDR = "0XABBAABBAABBAABBAABBAABBAABBAABBAABBAABBA"

var (
	// POAABI defines the ABI of the POA smart-contract as needed by a consensus
	// module to check if an address is authorized
	POAABI abi.ABI

	// POAABISTRING is the string representaion of POAABI
	POAABISTRING string

	// POAADDR is the address of the POA smart-contract.
	POAADDR eth_common.Address
)

func init() {
	POAABI, _ = abi.JSON(strings.NewReader(defaultPOAABI))
	POAADDR = eth_common.HexToAddress(defaultPOAADDR)
}

func setPOAABI(poaABI string) {
	POAABISTRING = poaABI
	POAABI, _ = abi.JSON(strings.NewReader(poaABI))
}

func setPOAADDR(address string) {
	POAADDR = eth_common.HexToAddress(address)
}
