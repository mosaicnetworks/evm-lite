package babble

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/mosaicnetworks/babble/crypto"
	"github.com/mosaicnetworks/babble/hashgraph"
	"github.com/mosaicnetworks/babble/net"
	"github.com/mosaicnetworks/babble/node"
	bserv "github.com/mosaicnetworks/babble/service"
	"github.com/mosaicnetworks/evm-lite/config"
	"github.com/mosaicnetworks/evm-lite/service"
	"github.com/mosaicnetworks/evm-lite/state"
	"github.com/sirupsen/logrus"
)

/*
InmemBabble implementes the Consensus Interface.
It uses an inmemory Babble node.
*/
type InmemBabble struct {
	config        config.BabbleConfig
	ethService    *service.Service
	ethState      *state.State
	babbleNode    *node.Node
	babbleService *bserv.Service
	logger        *logrus.Entry
}

//NewInmemBabble instantiates a new InmemBabble consensus system
func NewInmemBabble(config config.BabbleConfig, logger *logrus.Logger) *InmemBabble {
	return &InmemBabble{
		config: config,
		logger: logger.WithField("module", "babble"),
	}
}

/*******************************************************************************
IMPLEMENT CONSENSUS INTERFACE
*******************************************************************************/

//Init instantiates a Babble inmemory node
func (b *InmemBabble) Init(state *state.State, service *service.Service) error {

	b.logger.Debug("INIT")

	b.ethState = state
	b.ethService = service

	//--------------------------------------------------------------------------

	// Create the PEM key
	pemKey := crypto.NewPemKey(b.config.BabbleDir)

	// Try a read
	key, err := pemKey.ReadKey()
	if err != nil {
		return err
	}

	// Create the peer store
	peerStore := net.NewJSONPeers(b.config.BabbleDir)
	// Try a read
	peers, err := peerStore.Peers()
	if err != nil {
		return err
	}

	// There should be at least two peers
	if len(peers) < 2 {
		return fmt.Errorf("Should define at least two peers")
	}

	sort.Sort(net.ByPubKey(peers))
	pmap := make(map[string]int)
	for i, p := range peers {
		pmap[p.PubKeyHex] = i
	}

	//Find the ID of this node
	nodePub := fmt.Sprintf("0x%X", crypto.FromECDSAPub(&key.PublicKey))
	nodeID := pmap[nodePub]

	b.logger.WithFields(logrus.Fields{
		"pmap": pmap,
		"id":   nodeID,
	}).Debug("PARTICIPANTS")

	conf := node.NewConfig(
		time.Duration(b.config.Heartbeat)*time.Millisecond,
		time.Duration(b.config.TCPTimeout)*time.Millisecond,
		b.config.CacheSize,
		b.config.SyncLimit,
		b.config.StoreType,
		b.config.StorePath,
		logrus.New())

	//Instantiate the Store (inmem or badger)
	var store hashgraph.Store
	var needBootstrap bool
	switch conf.StoreType {
	case "inmem":
		store = hashgraph.NewInmemStore(pmap, conf.CacheSize)
	case "badger":
		//If the file already exists, load and bootstrap the store using the file
		if _, err := os.Stat(conf.StorePath); err == nil {
			b.logger.Debug("loading badger store from existing database")
			store, err = hashgraph.LoadBadgerStore(conf.CacheSize, conf.StorePath)
			if err != nil {
				return fmt.Errorf("failed to load BadgerStore from existing file: %s", err)
			}
			needBootstrap = true
		} else {
			//Otherwise create a new one
			b.logger.Debug("creating new badger store from fresh database")
			store, err = hashgraph.NewBadgerStore(pmap, conf.CacheSize, conf.StorePath)
			if err != nil {
				return fmt.Errorf("failed to create new BadgerStore: %s", err)
			}
		}
	default:
		return fmt.Errorf("Invalid StoreType: %s", conf.StoreType)
	}

	trans, err := net.NewTCPTransport(
		b.config.NodeAddr, nil, 2, conf.TCPTimeout, conf.Logger)
	if err != nil {
		return fmt.Errorf("Creating TCP Transport: %s", err)
	}

	appProxy := NewInmemProxy(state, service, service.GetSubmitCh(), b.logger)

	node := node.NewNode(conf, nodeID, key, peers, store, trans, appProxy)
	if err := node.Init(needBootstrap); err != nil {
		return fmt.Errorf("Initializing node: %s", err)
	}

	babbleService := bserv.NewService(b.config.BabbleAPIAddr, node, conf.Logger)

	//--------------------------------------------------------------------------
	b.babbleNode = node
	b.babbleService = babbleService
	//--------------------------------------------------------------------------

	return nil
}

//Run starts the Babble node and service
func (b *InmemBabble) Run() error {
	//Babble API service
	go b.babbleService.Serve()

	b.babbleNode.Run(true)

	return nil
}
