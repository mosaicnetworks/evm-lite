package config

import (
	"fmt"
)

var (
	defaultProxyAddr     = ":1339"
	defaultClientAddr    = ":1338"
	defaultNodeAddr      = "127.0.0.1:1337"
	defaultBabbleAPIAddr = ":8000"
	defaultHeartbeat     = 500
	defaultTCPTimeout    = 1000
	defaultCacheSize     = 50000
	defaultSyncLimit     = 1000
	defaultMaxPool       = 2
	defaultStoreType     = "badger"
	defaultBabbleDir     = fmt.Sprintf("%s/babble", DefaultDataDir)
	defaultPeersFile     = fmt.Sprintf("%s/peers.json", defaultBabbleDir)
	defaultStorePath     = fmt.Sprintf("%s/badger_db", defaultBabbleDir)
)

//BabbleConfig contains the configuration of a Babble node
type BabbleConfig struct {

	/*********************************************
	SOCKET
	*********************************************/

	//Address of Babble proxy
	ProxyAddr string `mapstructure:"proxy_addr"`

	//Address of Babble client proxy
	ClientAddr string `mapstructure:"client_addr"`

	/*********************************************
	Inmem
	*********************************************/

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
	StoreType string `mapstructure:"store_type"`

	//If StoreType = badger, location of database file
	StorePath string `mapstructure:"store_path"`
}

//DefaultBabbleConfig returns the default configuration for a Babble node
func DefaultBabbleConfig() *BabbleConfig {
	return &BabbleConfig{
		ProxyAddr:     defaultProxyAddr,
		ClientAddr:    defaultClientAddr,
		BabbleDir:     defaultBabbleDir,
		NodeAddr:      defaultNodeAddr,
		BabbleAPIAddr: defaultBabbleAPIAddr,
		Heartbeat:     defaultHeartbeat,
		TCPTimeout:    defaultTCPTimeout,
		CacheSize:     defaultCacheSize,
		SyncLimit:     defaultSyncLimit,
		MaxPool:       defaultMaxPool,
		StoreType:     defaultStoreType,
		StorePath:     defaultStorePath,
	}
}

//SetDataDir updates the babble configuration directories if they were set to
//to default values.
func (c *BabbleConfig) SetDataDir(datadir string) {
	if c.BabbleDir == defaultBabbleDir {
		c.BabbleDir = fmt.Sprintf("%s", datadir)
	}
	if c.StorePath == defaultStorePath {
		c.StorePath = fmt.Sprintf("%s/badger_db", c.BabbleDir)
	}
}
