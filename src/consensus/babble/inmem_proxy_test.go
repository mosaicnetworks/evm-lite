package babble

import (
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/mosaicnetworks/babble/src/hashgraph"
	"github.com/mosaicnetworks/babble/src/peers"
	"github.com/mosaicnetworks/evm-lite/src/state"
	"github.com/sirupsen/logrus"

	bcommon "github.com/mosaicnetworks/evm-lite/src/common"
)

var (
	_defaultValue    = big.NewInt(0)
	_defaultGas      = uint64(1000000)
	_defaultGasPrice = big.NewInt(0)
)

type Test struct {
	dataDir string
	pwdFile string
	dbFile  string
	cache   int

	keyStore *keystore.KeyStore
	state    *state.State
	logger   *logrus.Logger
}

func NewTest(dataDir string, logger *logrus.Logger, t *testing.T) *Test {
	pwdFile := filepath.Join(dataDir, "pwd.txt")
	dbFile := filepath.Join(dataDir, "chaindata")
	genesisFile := filepath.Join(dataDir, "genesis.json")
	cache := 128

	state, err := state.NewState(logger, dbFile, cache, genesisFile)
	if err != nil {
		t.Fatal(err)
	}

	return &Test{
		dataDir: dataDir,
		pwdFile: pwdFile,
		dbFile:  dbFile,
		cache:   cache,
		state:   state,
		logger:  logger,
	}
}

func (test *Test) readPwd() (pwd string, err error) {
	text, err := ioutil.ReadFile(test.pwdFile)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(text), "\n")
	// Sanitise DOS line endings.
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], "\r")
	}
	return lines[0], nil
}

func (test *Test) initKeyStore() error {
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP

	keydir := filepath.Join(test.dataDir, "keystore")
	if err := os.MkdirAll(keydir, 0700); err != nil {
		return err
	}

	test.keyStore = keystore.NewKeyStore(keydir, scryptN, scryptP)

	return nil
}

func (test *Test) unlockAccounts() error {
	pwd, err := test.readPwd()
	if err != nil {
		test.logger.WithError(err).Error("Reading PwdFile")
		return err
	}

	for _, ac := range test.keyStore.Accounts() {
		if err := test.keyStore.Unlock(ac, string(pwd)); err != nil {
			return err
		}
		test.logger.WithField("address", ac.Address.Hex()).Debug("Unlocked account")
	}
	return nil
}

func (test *Test) Init() error {
	if err := test.initKeyStore(); err != nil {
		return err
	}

	if err := test.unlockAccounts(); err != nil {
		return err
	}

	return nil
}

func (test *Test) prepareTransaction(from, to *accounts.Account,
	value *big.Int,
	gas uint64,
	gasPrice *big.Int,
	data []byte) (*ethTypes.Transaction, error) {

	nonce := test.state.GetPoolNonce(from.Address)

	var tx *ethTypes.Transaction
	if to == nil {
		tx = ethTypes.NewContractCreation(nonce,
			value,
			gas,
			gasPrice,
			data)
	} else {
		tx = ethTypes.NewTransaction(nonce,
			to.Address,
			value,
			gas,
			gasPrice,
			data)
	}

	signer := ethTypes.NewEIP155Signer(big.NewInt(1))

	signature, err := test.keyStore.SignHash(*from, signer.Hash(tx).Bytes())
	if err != nil {
		return nil, err
	}
	signedTx, err := tx.WithSignature(signer, signature)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func (test *Test) deployContract(from accounts.Account, contract *Contract, t *testing.T) {

	// Create Contract transaction
	tx, err := test.prepareTransaction(&from,
		nil,
		_defaultValue,
		_defaultGas,
		_defaultGasPrice,
		common.FromHex(contract.code))

	if err != nil {
		t.Fatal(err)
	}

	// Convert to raw bytes
	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		t.Fatal(err)
	}

	// Try to commit the transaction
	err = test.state.ApplyTransaction(data, 0, common.Hash{})
	if err != nil {
		t.Fatal(err)
	}
	_, err = test.state.Commit()
	if err != nil {
		t.Fatal(err)
	}

	receipt, err := test.state.GetReceipt(tx.Hash())
	if err != nil {
		t.Fatal(err)
	}

	contract.address = receipt.ContractAddress
}

type Contract struct {
	name    string
	address common.Address
	code    string
	abi     string
	jsonABI abi.ABI
}

func (c *Contract) parseABI(t *testing.T) {
	jABI, err := abi.JSON(strings.NewReader(c.abi))
	if err != nil {
		t.Fatal(err)
	}
	c.jsonABI = jABI
}

//------------------------------------------------------------------------------
// Always return true test

/*

pragma solidity >=0.4.0;

contract Test {

  function checkAuthorisedPublicKey(bytes32  _publicKey) constant returns (bool) {

      return true;
   }
}

*/

func alwaysTrueContract() *Contract {
	return &Contract{
		name: "Test",
		code: "606060405234610000575b60aa806100186000396000f30060606040526000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630eefa3ab14603c575b6000565b3460005760586004808035600019169060200190919050506072565b604051808215151515815260200191505060405180910390f35b6000600190505b9190505600a165627a7a72305820bdf53f15fe9b0bfc0f922b6b8eec73d74eb30969090231af9e882e3f131e7a460029",
		abi:  "[{\"constant\": true,\"inputs\": [{\"name\": \"_publicKey\",\"type\": \"bytes32\"}],\"name\": \"checkAuthorisedPublicKey\",\"outputs\": [{\"name\": \"\",\"type\": \"bool\"}],\"payable\": false,\"type\": \"function\",\"stateMutability\": \"view\"}]",
	}
}

func TestAlwaysTrueContract(t *testing.T) {

	os.RemoveAll("test_data/eth/chaindata")
	defer os.RemoveAll("test_data/eth/chaindata")

	testLogger := bcommon.NewTestLogger(t)

	test := NewTest("test_data/eth", testLogger, t)
	//defer test.state.db.Close()

	err := test.Init()

	if err != nil {
		t.Fatal(err)
	}

	from := test.keyStore.Accounts()[0]

	contract := alwaysTrueContract()

	test.deployContract(from, contract, t)

	contract.parseABI(t)

	inmemProxy := &InmemProxy{
		state:  test.state,
		logger: testLogger.WithField("module", "babble/proxy"),
	}

	peerSlice := []*peers.Peer{peers.NewPeer("0x1234", "0.0.0.0", "test")}
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

	if res.InternalTransactions[0].Accepted != hashgraph.True {
		t.Fatalf("Result wrong")
	}

}

//------------------------------------------------------------------------------
// Mixed test

/*

pragma solidity 0.5.7;

contract Test {

  function checkAuthorisedPublicKey(bytes32  _publicKey) public view returns (bool) {

      if(_publicKey == "0x12345") {
          return true;
      } else {
          return false;
      }
   }
}

*/

//0.4.8     y
//0.5.8     n
//0.5.6     n
//0.4.26    y
//0.5.0     y
//0.5.7     n
//0.5.4     y
//0.5.5     n

func mixedContract() *Contract {
	return &Contract{
		name: "Test",
		//code: "606060405234610000575b60e1806100186000396000f30060606040526000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630eefa3ab14603c575b6000565b3460005760586004808035600019169060200190919050506072565b604051808215151515815260200191505060405180910390f35b60007f30783132333435000000000000000000000000000000000000000000000000008260001916141560a7576001905060b0565b6001905060b0565b5b9190505600a165627a7a723058205385d98f265d3a6043a701657984770fe1fa22658f357a673bc348b6d4bd054d0029",
		//code: "608060405234801561001057600080fd5b5060d78061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80630eefa3ab14602d575b600080fd5b605660048036036020811015604157600080fd5b81019080803590602001909291905050506070565b604051808215151515815260200191505060405180910390f35b60007f307831323334350000000000000000000000000000000000000000000000000082141560a1576001905060a6565b600190505b91905056fea165627a7a7230582084b52999fe6cae70801f1f59e13f60782ad61ba5cd20003a9430f888733a723c0029",
		//code: "608060405234801561001057600080fd5b5060d78061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80630eefa3ab14602d575b600080fd5b605660048036036020811015604157600080fd5b81019080803590602001909291905050506070565b604051808215151515815260200191505060405180910390f35b60007f307831323334350000000000000000000000000000000000000000000000000082141560a1576001905060a6565b600190505b91905056fea165627a7a723058204ffbf4b3e536e972926a1a883415870176c02a3b182693c3428a679b6a9dda460029",
		//code: "608060405234801561001057600080fd5b5060f58061001f6000396000f300608060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630eefa3ab146044575b600080fd5b348015604f57600080fd5b5060706004803603810190808035600019169060200190929190505050608a565b604051808215151515815260200191505060405180910390f35b60007f30783132333435000000000000000000000000000000000000000000000000008260001916141560bf576001905060c4565b600190505b9190505600a165627a7a72305820035ef062221c725736f13a9576fefc58220fea4fdace2b213073a9f92be1ca6b0029",
		//code: "608060405234801561001057600080fd5b5060fa8061001f6000396000f3fe608060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630eefa3ab146044575b600080fd5b348015604f57600080fd5b50607960048036036020811015606457600080fd5b81019080803590602001909291905050506093565b604051808215151515815260200191505060405180910390f35b60007f307831323334350000000000000000000000000000000000000000000000000082141560c4576001905060c9565b600190505b91905056fea165627a7a72305820ff180acc0b1bca66eb43e15db50a59db93c492e24feaedbc11e1f4c8d530c6930029",
		//code: "608060405234801561001057600080fd5b5060d78061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80630eefa3ab14602d575b600080fd5b605660048036036020811015604157600080fd5b81019080803590602001909291905050506070565b604051808215151515815260200191505060405180910390f35b60007f307831323334350000000000000000000000000000000000000000000000000082141560a1576001905060a6565b600190505b91905056fea165627a7a72305820a7ce6e7b31e11070466e2a29054735af94afb286b3b7fcbed8b40a8885ed17ba0029",
		//code: "608060405234801561001057600080fd5b5060f48061001f6000396000f3fe6080604052348015600f57600080fd5b50600436106045576000357c0100000000000000000000000000000000000000000000000000000000900480630eefa3ab14604a575b600080fd5b607360048036036020811015605e57600080fd5b8101908080359060200190929190505050608d565b604051808215151515815260200191505060405180910390f35b60007f307831323334350000000000000000000000000000000000000000000000000082141560be576001905060c3565b600190505b91905056fea165627a7a72305820c2571c791e66c7211162e58c716aac447fecc528ba227a592aed238bf50110660029",
		//code: "608060405234801561001057600080fd5b5060d78061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80630eefa3ab14602d575b600080fd5b605660048036036020811015604157600080fd5b81019080803590602001909291905050506070565b604051808215151515815260200191505060405180910390f35b60007f307831323334350000000000000000000000000000000000000000000000000082141560a1576001905060a6565b600190505b91905056fea165627a7a72305820f34a607615ed048aeb9772302d95f7288b5331fea28674084d9b6af31e6008630029",
		code: "608060405234801561001057600080fd5b5060d78061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80630eefa3ab14602d575b600080fd5b605660048036036020811015604157600080fd5b81019080803590602001909291905050506070565b604051808215151515815260200191505060405180910390f35b60007f307831323334350000000000000000000000000000000000000000000000000082141560a1576001905060a6565b600090505b91905056fea165627a7a72305820be45bfd0f46cdc7089f5a5cbea1cb1e6ca9c872526a739ada54c1fa0d82614c10029",
		abi:  "[{\"constant\": true,\"inputs\": [{\"name\": \"_publicKey\",\"type\": \"bytes32\"}],\"name\": \"checkAuthorisedPublicKey\",\"outputs\": [{\"name\": \"\",\"type\": \"bool\"}],\"payable\": false,\"type\": \"function\",\"stateMutability\": \"view\"}]",
	}
}

func TestMixedContract(t *testing.T) {

	os.RemoveAll("test_data/eth/chaindata")
	defer os.RemoveAll("test_data/eth/chaindata")

	testLogger := bcommon.NewTestLogger(t)

	test := NewTest("test_data/eth", testLogger, t)
	//defer test.state.db.Close()

	err := test.Init()

	if err != nil {
		t.Fatal(err)
	}

	from := test.keyStore.Accounts()[0]

	contract := mixedContract()

	test.deployContract(from, contract, t)

	contract.parseABI(t)

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

	if res.InternalTransactions[0].Accepted != hashgraph.True {
		t.Log(contract.address.String())
		t.Log(res.InternalTransactions[0].Accepted)
		t.Fatalf("Result wrong")
	}

}
