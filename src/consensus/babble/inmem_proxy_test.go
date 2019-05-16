package babble

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mosaicnetworks/babble/src/common"
	"github.com/mosaicnetworks/babble/src/hashgraph"
	"github.com/mosaicnetworks/babble/src/peers"
	"github.com/mosaicnetworks/evm-lite/src/state"
	"github.com/sirupsen/logrus"

	bcommon "github.com/mosaicnetworks/evm-lite/src/common"
)

// pragma solidity 0.5.7;

// contract Test {

//   function checkAuthorisedPublicKey(bytes32 _publicKey) public pure returns (bool) {

//       if(_publicKey == "0x12345") {
//           return true;
//       } else {
//           return false;
//       }
//    }
// }

// solc --bin-runtime --overwrite -o out

type Test struct {
	dataDir string
	dbFile  string
	cache   int

	state  *state.State
	logger *logrus.Logger
}

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

	peerSlice := []*peers.Peer{peers.NewPeer("0x12345", "0.0.0.0", "test")}
	var txs [][]byte
	itxs := []hashgraph.InternalTransaction{hashgraph.NewInternalTransaction(hashgraph.PEER_ADD, *peerSlice[0])}

	block := hashgraph.NewBlock(1,
		1,
		[]byte("frameHash"),
		peerSlice,
		txs,
		itxs)

	res, _ := inmemProxy.CommitBlock(*block)

	if len(res.InternalTransactions) < 1 {
		t.Fatalf("Length response too short")
	}

	if res.InternalTransactions[0].Accepted != common.True {
		t.Log(res.InternalTransactions[0].Accepted)
		t.Fatalf("Result wrong")
	}

}
