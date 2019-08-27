package run

import (
	_config "github.com/mosaicnetworks/evm-lite/src/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config = _config.DefaultConfig()
	logger = logrus.New()
)

//RunCmd is launches a node
var RunCmd = &cobra.Command{
	Use:              "run",
	Short:            "Run a node",
	TraverseChildren: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err := bindFlagsLoadViper(cmd); err != nil {
			return err
		}

		config, err = parseConfig()
		if err != nil {
			return err
		}

		logger = logrus.New()
		logger.Level = logLevel(config.LogLevel)

		config.SetDataDir(config.DataDir)

		logger.WithFields(logrus.Fields{
			"Base": config}).Debug("Config")

		return nil
	},
}

func init() {
	//Subcommands
	RunCmd.AddCommand(
		NewSoloCmd())

	//Base config
	RunCmd.PersistentFlags().StringP("datadir", "d", config.DataDir, "Top-level directory for configuration and data")
	RunCmd.PersistentFlags().String("log", config.LogLevel, "debug, info, warn, error, fatal, panic")

	//Eth config
	RunCmd.PersistentFlags().String("eth.genesis", config.Genesis, "Location of genesis file")
	RunCmd.PersistentFlags().String("eth.keystore", config.Keystore, "Location of Ethereum account keys")
	RunCmd.PersistentFlags().String("eth.pwd", config.PwdFile, "Password file to unlock accounts")
	RunCmd.PersistentFlags().String("eth.db", config.DbFile, "Eth database file")
	RunCmd.PersistentFlags().String("eth.listen", config.EthAPIAddr, "Address of HTTP API service")
	RunCmd.PersistentFlags().Int("eth.cache", config.Cache, "Megabytes of memory allocated to internal caching (min 16MB / database forced)")

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

	viper.SetConfigName("evml")         // name of config file (without extension)
	viper.AddConfigPath(config.DataDir) // search root directory

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// stderr, so if we redirect output to json file, this doesn't appear
		logger.Debugf("Using config file: %s", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		logger.Debugf("No config file found in %s", config.DataDir)
	} else {

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
