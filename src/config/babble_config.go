package config

import (
	"fmt"
	"time"

	_babble "github.com/mosaicnetworks/babble/src/babble"
	"github.com/sirupsen/logrus"
)

var (
	defaultNodeAddr      = "127.0.0.1:1337"
	defaultBabbleAPIAddr = ":8000"
	defaultHeartbeat     = 500
	defaultTCPTimeout    = 1000
	defaultCacheSize     = 50000
	defaultSyncLimit     = 1000
	defaultMaxPool       = 2
	defaultBabbleDir     = fmt.Sprintf("%s/babble", DefaultDataDir)
	defaultPeersFile     = fmt.Sprintf("%s/peers.json", defaultBabbleDir)
)

//BabbleConfig contains the configuration of a Babble node
type BabbleConfig struct {

	//Directory containing priv_key.pem and peers.json files
	BabbleDir string `mapstructure:"dir"`

	//Address of Babble node (where it talks to other Babble nodes)
	NodeAddr string `mapstructure:"node_addr"`

	//Babble HTTP API address
	BabbleAPIAddr string `mapstructure:"api_addr"`

	//Gossip heartbeat in milliseconds
	Heartbeat int `mapstructure:"heartbeat"`

	//TCP timeout in milliseconds
	TCPTimeout int `mapstructure:"tcp_timeout"`

	//Max number of items in caches
	CacheSize int `mapstructure:"cache_size"`

	//Max number of Event in SyncResponse
	SyncLimit int `mapstructure:"sync_limit"`

	//Max number of connections in net pool
	MaxPool int `mapstructure:"max_pool"`

	//Database type; badger or inmeum
	Store bool `mapstructure:"store"`
}

//DefaultBabbleConfig returns the default configuration for a Babble node
func DefaultBabbleConfig() *BabbleConfig {
	return &BabbleConfig{
		BabbleDir:     defaultBabbleDir,
		NodeAddr:      defaultNodeAddr,
		BabbleAPIAddr: defaultBabbleAPIAddr,
		Heartbeat:     defaultHeartbeat,
		TCPTimeout:    defaultTCPTimeout,
		CacheSize:     defaultCacheSize,
		SyncLimit:     defaultSyncLimit,
		MaxPool:       defaultMaxPool,
	}
}

//SetDataDir updates the babble configuration directories if they were set to
//to default values.
func (c *BabbleConfig) SetDataDir(datadir string) {
	if c.BabbleDir == defaultBabbleDir {
		c.BabbleDir = fmt.Sprintf("%s", datadir)
	}
}

//ToRealBabbleConfig converts an evm-lite/src/config.BabbleConfig to a
//babble/src/babble.BabbleConfig as used by Babble
func (c *BabbleConfig) ToRealBabbleConfig(logger *logrus.Logger) *_babble.BabbleConfig {
	babbleConfig := _babble.NewDefaultConfig()
	babbleConfig.DataDir = c.BabbleDir
	babbleConfig.BindAddr = c.NodeAddr
	babbleConfig.MaxPool = c.MaxPool
	babbleConfig.Store = c.Store
	babbleConfig.Logger = logger
	babbleConfig.NodeConfig.HeartbeatTimeout = time.Duration(c.Heartbeat) * time.Millisecond
	babbleConfig.NodeConfig.TCPTimeout = time.Duration(c.TCPTimeout) * time.Millisecond
	babbleConfig.NodeConfig.CacheSize = c.CacheSize
	babbleConfig.NodeConfig.SyncLimit = c.SyncLimit
	babbleConfig.NodeConfig.Logger = logger
	return babbleConfig
}
