package commands

import (
	"path/filepath"

	_config "github.com/mosaicnetworks/evm-lite/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config = _config.DefaultConfig()
	logger = logrus.New()
)

//RootCmd is the root command for evm-babble
var RootCmd = &cobra.Command{
	Use:              "evm-lite",
	Short:            "LightWeight EVM app for different consensus sytems",
	TraverseChildren: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if cmd.Name() == VersionCmd.Name() {
			return nil
		}

		if err := bindFlagsLoadViper(cmd); err != nil {
			return err
		}

		config, err = parseConfig()
		if err != nil {
			return err
		}

		logger = logrus.New()
		logger.Level = logLevel(config.BaseConfig.LogLevel)

		logger.WithFields(logrus.Fields{
			"Base": config.BaseConfig,
			"Eth":  config.Eth}).Debug("Config")

		return nil
	},
}

func init() {
	//Base
	RootCmd.PersistentFlags().String("datadir", config.BaseConfig.DataDir, "Top-level directory for configuration and data")
	RootCmd.PersistentFlags().String("log_level", config.BaseConfig.LogLevel, "debug, info, warn, error, fatal, panic")

	//Eth
	RootCmd.PersistentFlags().String("eth.genesis", config.Eth.Genesis, "Location of genesis file")
	RootCmd.PersistentFlags().String("eth.keystore", config.Eth.Keystore, "Location of Ethereum account keys")
	RootCmd.PersistentFlags().String("eth.pwd", config.Eth.PwdFile, "Password file to unlock accounts")
	RootCmd.PersistentFlags().String("eth.db", config.Eth.DbFile, "Eth database file")
	RootCmd.PersistentFlags().String("eth.api_addr", config.Eth.EthAPIAddr, "Address of HTTP API service")
	RootCmd.PersistentFlags().Int("eth.cache", config.Eth.Cache, "Megabytes of memory allocated to internal caching (min 16MB / database forced)")

}

//------------------------------------------------------------------------------

//Retrieve the default environment configuration.
func parseConfig() (*_config.Config, error) {
	conf := _config.DefaultConfig()
	err := viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}
	return conf, err
}

//Bind all flags and read the config into viper
func bindFlagsLoadViper(cmd *cobra.Command) error {
	// cmd.Flags() includes flags from this command and all persistent flags from the parent
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	viper.SetConfigName("config")                                           // name of config file (without extension)
	viper.AddConfigPath(config.BaseConfig.DataDir)                          // search root directory
	viper.AddConfigPath(filepath.Join(config.BaseConfig.DataDir, "config")) // search root directory /config

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// stderr, so if we redirect output to json file, this doesn't appear
		// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		// ignore not found error, return other errors
		return err
	}

	return nil
}

func logLevel(l string) logrus.Level {
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
