package babble

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mosaicnetworks/babble/src/hashgraph"
	"github.com/mosaicnetworks/babble/src/peers"
	"github.com/mosaicnetworks/evm-lite/src/state"
	"github.com/sirupsen/logrus"

	bcommon "github.com/mosaicnetworks/evm-lite/src/common"
)

/*

The goal of this test is to check the link between InternalTransactions and the
POA smart-contract.

An InternalTransaction contains the public key of the peer who is requesting to
join. The inmem-poxy needs to check in the smart-contract if the address
corresponding to the public-key is authorised.

Note that the inmem_proxy has to convert the public-key into an address before
calling the POA smart-contract.

The smart-contract needs to expose a checkAuthorised(address) method that
returns true or false. That is the only requirement at this stage. So for this
test, we use a dummy POA contract defined as follows:

	pragma solidity 0.5.7;

	contract Test {
		function checkAuthorised(address _address) public pure returns (bool) {
			if(_address == address(0x89acCD6b63d6eE73550eca0Cba16C2027c13FDa6)) {
			return true;
			} else {
			return false;
			}
		}
	}

It basicaly checks if the input address matches the following address:

0x89acCD6b63d6eE73550eca0Cba16C2027c13FDa6

which corresponds to the following public key:

0x04a9570b06f6e815d5b1e74eb30e7a7487d6589d9095884daf62d9b04f6542de39bdbb68b7cc41957ad15d699c98baf7fa18b12e638aa33e215d70c9aafd6c6c1d

The genesis.json file in the test_data/ directory defines the smart-contract
account with the corresponding bytecode, which can be obtained with the online
remix compiler or solc (solc --bin-runtime --overwrite -o out).

*/

type Test struct {
	dataDir string
	dbFile  string
	cache   int

	state  *state.State
	logger *logrus.Logger
}

/*
Pubkey : 0x04a9570b06f6e815d5b1e74eb30e7a7487d6589d9095884daf62d9b04f6542de39bdbb68b7cc41957ad15d699c98baf7fa18b12e638aa33e215d70c9aafd6c6c1d
Address: 0x89acCD6b63d6eE73550eca0Cba16C2027c13FDa6
*/
const authPubkey = "0x04a9570b06f6e815d5b1e74eb30e7a7487d6589d9095884daf62d9b04f6542de39bdbb68b7cc41957ad15d699c98baf7fa18b12e638aa33e215d70c9aafd6c6c1d"

func NewTest(dataDir string, logger *logrus.Logger, t *testing.T) *Test {
	dbFile := filepath.Join(dataDir, "chaindata")
	genesisFile := filepath.Join(dataDir, "genesis.json")
	cache := 128

	state, err := state.NewState(logger, dbFile, cache, genesisFile)
	if err != nil {
		t.Fatal(err)
	}

	return &Test{
		dataDir: dataDir,
		dbFile:  dbFile,
		cache:   cache,
		state:   state,
		logger:  logger,
	}
}

func TestMixedContract(t *testing.T) {

	os.RemoveAll("test_data/eth/chaindata")
	defer os.RemoveAll("test_data/eth/chaindata")

	testLogger := bcommon.NewTestLogger(t)

	test := NewTest("test_data/eth", testLogger, t)
	//defer test.state.db.Close()

	inmemProxy := &InmemProxy{
		state:  test.state,
		logger: testLogger.WithField("module", "babble/proxy"),
	}

	peerSlice := []*peers.Peer{peers.NewPeer(authPubkey, "0.0.0.0", "test")}
	var txs [][]byte
	itxs := []hashgraph.InternalTransaction{hashgraph.NewInternalTransaction(hashgraph.PEER_ADD, *peerSlice[0])}

	block := hashgraph.NewBlock(1,
		1,
		[]byte("frameHash"),
		peerSlice,
		txs,
		itxs)

	res, _ := inmemProxy.CommitBlock(*block)

	if len(res.InternalTransactionReceipts) < 1 {
		t.Fatal("Length response too short")
	}

	if !res.InternalTransactionReceipts[0].Accepted {
		t.Fatal("InternalTransactionReceipts[0] should be accepted")
	}

}
