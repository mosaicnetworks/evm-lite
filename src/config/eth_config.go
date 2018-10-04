package config

import "fmt"

var (
	defaultEthAPIAddr   = ":8080"
	defaultCache        = 128
	defaultEthDir       = fmt.Sprintf("%s/eth", DefaultDataDir)
	defaultKeystoreFile = fmt.Sprintf("%s/keystore", defaultEthDir)
	defaultGenesisFile  = fmt.Sprintf("%s/genesis.json", defaultEthDir)
	defaultPwdFile      = fmt.Sprintf("%s/pwd.txt", defaultEthDir)
	defaultDbFile       = fmt.Sprintf("%s/chaindata", defaultEthDir)
)

//EthConfig contains the configuration relative to the accounts, EVM, trie/db,
//and service API
type EthConfig struct {

	//Genesis file
	Genesis string `mapstructure:"genesis"`

	//Location of ethereum account keys
	Keystore string `mapstructure:"keystore"`

	//File containing passwords to unlock ethereum accounts
	PwdFile string `mapstructure:"pwd"`

	//File containing the levelDB database
	DbFile string `mapstructure:"db"`

	//Address of HTTP API Service
	EthAPIAddr string `mapstructure:"listen"`

	//Megabytes of memory allocated to internal caching (min 16MB / database forced)
	Cache int `mapstructure:"cache"`
}

//DefaultEthConfig return the default configuration for Eth services
func DefaultEthConfig() *EthConfig {
	return &EthConfig{
		Genesis:    defaultGenesisFile,
		Keystore:   defaultKeystoreFile,
		PwdFile:    defaultPwdFile,
		DbFile:     defaultDbFile,
		EthAPIAddr: defaultEthAPIAddr,
		Cache:      defaultCache,
	}
}

//SetDataDir updates the eth configuration directories if they were set to
//default values.
func (c *EthConfig) SetDataDir(datadir string) {
	if c.Genesis == defaultGenesisFile {
		c.Genesis = fmt.Sprintf("%s/genesis.json", datadir)
	}
	if c.Keystore == defaultKeystoreFile {
		c.Keystore = fmt.Sprintf("%s/keystore", datadir)
	}
	if c.PwdFile == defaultPwdFile {
		c.PwdFile = fmt.Sprintf("%s/pwd.txt", datadir)
	}
	if c.DbFile == defaultDbFile {
		c.DbFile = fmt.Sprintf("%s/chaindata", datadir)
	}
}
