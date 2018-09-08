package engine

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

var (
	//Base
	defaultLogLevel = "debug"
	defaultDataDir  = defaultHomeDir()

	//Eth
	defaultEthAPIAddr   = ":8080"
	defaultCache        = 128
	defaultEthDir       = fmt.Sprintf("%s/eth", defaultDataDir)
	defaultKeystoreFile = fmt.Sprintf("%s/keystore", defaultEthDir)
	defaultGenesisFile  = fmt.Sprintf("%s/genesis.json", defaultEthDir)
	defaultPwdFile      = fmt.Sprintf("%s/pwd.txt", defaultEthDir)
	defaultDbFile       = fmt.Sprintf("%s/chaindata", defaultEthDir)
)

//Config contains de configuration for an EVM-Lite node
type Config struct {

	//Top level options use an anonymous struct
	BaseConfig `mapstructure:",squash"`

	//Options for EVM and State
	Eth *EthConfig `mapstructure:"eth"`
}

//DefaultConfig returns the default configuration for an EVM-Lite node
func DefaultConfig() *Config {
	return &Config{
		BaseConfig: DefaultBaseConfig(),
		Eth:        DefaultEthConfig(),
	}
}

/*******************************************************************************
BASE CONFIG
*******************************************************************************/

//BaseConfig contains the top level configuration for an EVM-Babble node
type BaseConfig struct {

	//Top-level directory of evm-babble data
	DataDir string `mapstructure:"datadir"`

	//Debug, info, warn, error, fatal, panic
	LogLevel string `mapstructure:"log_level"`
}

//DefaultBaseConfig returns the default top-level configuration for EVM-Babble
func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		DataDir:  defaultDataDir,
		LogLevel: defaultLogLevel,
	}
}

/*******************************************************************************
ETH CONFIG
*******************************************************************************/

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
	EthAPIAddr string `mapstructure:"api_addr"`

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

/*******************************************************************************
FILE HELPERS
*******************************************************************************/

func defaultHomeDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "BABBLE")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "EVMBABBE")
		} else {
			return filepath.Join(home, ".evm-lite")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}
