package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	// Base
	defaultLogLevel     = "debug"
	defaultDataDir      = defaultHomeDir()
	defaultEthAPIAddr   = ":8080"
	defaultCache        = 128
	defaultEthDir       = fmt.Sprintf("%s/eth", defaultDataDir)
	defaultKeystoreFile = fmt.Sprintf("%s/keystore", defaultEthDir)
	defaultGenesisFile  = fmt.Sprintf("%s/genesis.json", defaultEthDir)
	defaultPwdFile      = fmt.Sprintf("%s/pwd.txt", defaultEthDir)
	defaultDbFile       = fmt.Sprintf("%s/chaindata", defaultEthDir)
)

// Config contains de configuration for an EVM-Lite node
type Config struct {
	// Top-level directory of evm-babble data
	DataDir string `mapstructure:"datadir"`

	// Debug, info, warn, error, fatal, panic
	LogLevel string `mapstructure:"log"`

	// Genesis file
	Genesis string `mapstructure:"genesis"`

	// Location of ethereum account keys
	Keystore string `mapstructure:"keystore"`

	// File containing passwords to unlock ethereum accounts
	PwdFile string `mapstructure:"pwd"`

	// File containing the levelDB database
	DbFile string `mapstructure:"db"`

	// Address of HTTP API Service
	EthAPIAddr string `mapstructure:"listen"`

	// Megabytes of memory allocated to internal caching (min 16MB / database
	// forced)
	Cache int `mapstructure:"cache"`

	logger *logrus.Logger
}

// DefaultConfig returns the default configuration for an EVM-Lite node
func DefaultConfig() *Config {
	return &Config{
		DataDir:    defaultDataDir,
		LogLevel:   defaultLogLevel,
		Genesis:    defaultGenesisFile,
		Keystore:   defaultKeystoreFile,
		PwdFile:    defaultPwdFile,
		DbFile:     defaultDbFile,
		EthAPIAddr: defaultEthAPIAddr,
		Cache:      defaultCache,
	}
}

// SetDataDir updates the root data directory and trickles down to the eth
// directories if they are currently set to the default values.
func (c *Config) SetDataDir(datadir string) {
	c.DataDir = datadir

	if c.Genesis == defaultGenesisFile {
		c.Genesis = fmt.Sprintf("%s/eth/genesis.json", datadir)
	}
	if c.Keystore == defaultKeystoreFile {
		c.Keystore = fmt.Sprintf("%s/eth/keystore", datadir)
	}
	if c.PwdFile == defaultPwdFile {
		c.PwdFile = fmt.Sprintf("%s/eth/pwd.txt", datadir)
	}
	if c.DbFile == defaultDbFile {
		c.DbFile = fmt.Sprintf("%s/eth/chaindata", datadir)
	}
}

// Logger returns a formatted logrus Entry that supports nested prefixes.
func (c *Config) Logger() *logrus.Entry {
	if c.logger == nil {
		c.logger = logrus.New()
		c.logger.Level = LogLevel(c.LogLevel)
		c.logger.Formatter = new(prefixed.TextFormatter)
	}
	return c.logger.WithField("prefix", "evm-lite")
}

// LogLevel ...
func LogLevel(l string) logrus.Level {
	switch l {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.DebugLevel
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
			return filepath.Join(home, "Library", "EVMLITE")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "EVMLITE")
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
