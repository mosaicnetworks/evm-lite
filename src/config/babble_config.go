package config

import (
	"fmt"
	"time"

	_babble "github.com/mosaicnetworks/babble/src/babble"
)

var (
	defaultNodeAddr       = ":1337"
	defaultBabbleAPIAddr  = ":8000"
	defaultHeartbeat      = 500 * time.Millisecond
	defaultTCPTimeout     = 1000 * time.Millisecond
	defaultCacheSize      = 50000
	defaultSyncLimit      = 1000
	defaultEnableFastSync = true
	defaultMaxPool        = 2
	defaultBabbleDir      = fmt.Sprintf("%s/babble", DefaultDataDir)
	defaultPeersFile      = fmt.Sprintf("%s/peers.json", defaultBabbleDir)
)

// BabbleConfig contains the configuration of a Babble node
type BabbleConfig struct {

	// Directory containing priv_key.pem and peers.json files
	DataDir string `mapstructure:"datadir"`

	// Address of Babble node (where it talks to other Babble nodes)
	BindAddr string `mapstructure:"listen"`

	// Babble HTTP API address
	ServiceAddr string `mapstructure:"service-listen"`

	// Gossip heartbeat
	Heartbeat time.Duration `mapstructure:"heartbeat"`

	// TCP timeout
	TCPTimeout time.Duration `mapstructure:"timeout"`

	// Max number of items in caches
	CacheSize int `mapstructure:"cache-size"`

	// Max number of Event in SyncResponse
	SyncLimit int `mapstructure:"sync-limit"`

	// Allow node to FastSync
	EnableFastSync bool `mapstructure:"enable-fast-sync"`

	// Max number of connections in net pool
	MaxPool int `mapstructure:"max-pool"`

	// Database type; badger or inmeum
	Store bool `mapstructure:"store"`
}

// DefaultBabbleConfig returns the default configuration for a Babble node
func DefaultBabbleConfig() *BabbleConfig {
	return &BabbleConfig{
		DataDir:        defaultBabbleDir,
		BindAddr:       defaultNodeAddr,
		ServiceAddr:    defaultBabbleAPIAddr,
		Heartbeat:      defaultHeartbeat,
		TCPTimeout:     defaultTCPTimeout,
		CacheSize:      defaultCacheSize,
		SyncLimit:      defaultSyncLimit,
		EnableFastSync: defaultEnableFastSync,
		MaxPool:        defaultMaxPool,
	}
}

// SetDataDir updates the babble configuration directories if they were set to
// to default values.
func (c *BabbleConfig) SetDataDir(datadir string) {
	if c.DataDir == defaultBabbleDir {
		c.DataDir = datadir
	}
}

// ToRealBabbleConfig converts an evm-lite/src/config.BabbleConfig to a
// babble/src/babble.BabbleConfig as used by Babble
func (c *BabbleConfig) ToRealBabbleConfig() *_babble.BabbleConfig {
	babbleConfig := _babble.NewDefaultConfig()
	babbleConfig.DataDir = c.DataDir
	babbleConfig.BindAddr = c.BindAddr
	babbleConfig.ServiceAddr = c.ServiceAddr
	babbleConfig.MaxPool = c.MaxPool
	babbleConfig.Store = c.Store
	babbleConfig.NodeConfig.HeartbeatTimeout = c.Heartbeat
	babbleConfig.NodeConfig.TCPTimeout = c.TCPTimeout
	babbleConfig.NodeConfig.CacheSize = c.CacheSize
	babbleConfig.NodeConfig.SyncLimit = c.SyncLimit
	babbleConfig.NodeConfig.EnableFastSync = c.EnableFastSync
	return babbleConfig
}
