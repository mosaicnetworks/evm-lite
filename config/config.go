package config

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

var (
	//Base
	defaultLogLevel = "debug"
	DefaultDataDir  = defaultHomeDir()
)

//Config contains de configuration for an EVM-Lite node
type Config struct {

	//Top level options use an anonymous struct
	BaseConfig `mapstructure:",squash"`

	//Options for EVM and State
	Eth *EthConfig `mapstructure:"eth"`

	//Options for Babble consensus
	Babble *BabbleConfig `mapstructure:"babble"`
}

//DefaultConfig returns the default configuration for an EVM-Lite node
func DefaultConfig() *Config {
	return &Config{
		BaseConfig: DefaultBaseConfig(),
		Eth:        DefaultEthConfig(),
		Babble:     DefaultBabbleConfig(),
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
		DataDir:  DefaultDataDir,
		LogLevel: defaultLogLevel,
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
