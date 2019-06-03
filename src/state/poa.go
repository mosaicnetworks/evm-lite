package state

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	eth_common "github.com/ethereum/go-ethereum/common"
)

const poaABI = "[{\"type\":\"function\",\"inputs\": [{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"checkAuthorised\",\"outputs\": [{\"name\":\"\",\"type\":\"bool\"}]}]"
const poaFrom = "0x1337133713371337133713371337133713371337"

var (
	// POAABI defines the ABI of the POA smart-contract as needed by a consensus
	// module to check if an address is authorized
	POAABI abi.ABI

	// POAFROM is the address used in the 'from' field when querying the POA
	// smart-contract.
	POAFROM eth_common.Address
)

func init() {
	POAABI, _ = abi.JSON(strings.NewReader(poaABI))
	POAFROM = eth_common.HexToAddress(poaFrom)
}
