package commands

import (
	"fmt"

	"github.com/mosaicnetworks/evm-lite/consensus/babble"
	"github.com/mosaicnetworks/evm-lite/engine"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//AddBabbleFlags adds flags to the Babble command
func AddBabbleFlags(cmd *cobra.Command) {

	cmd.Flags().String("babble.proxy_addr", config.Babble.ProxyAddr, "IP:PORT of Babble proxy")
	cmd.Flags().String("babble.client_addr", config.Babble.ClientAddr, "IP:PORT to bind client proxy")
	cmd.Flags().String("babble.dir", config.Babble.BabbleDir, "Directory contaning priv_key.pem and peers.json files")
	cmd.Flags().String("babble.node_addr", config.Babble.NodeAddr, "IP:PORT of Babble node")
	cmd.Flags().String("babble.api_addr", config.Babble.BabbleAPIAddr, "IP:PORT of Babble HTTP API service")
	cmd.Flags().Int("babble.heartbeat", config.Babble.Heartbeat, "Heartbeat time milliseconds (time between gossips)")
	cmd.Flags().Int("babble.tcp_timeout", config.Babble.TCPTimeout, "TCP timeout milliseconds")
	cmd.Flags().Int("babble.cache_size", config.Babble.CacheSize, "Number of items in LRU caches")
	cmd.Flags().Int("babble.sync_limit", config.Babble.SyncLimit, "Max number of Events per sync")
	cmd.Flags().Int("babble.max_pool", config.Babble.MaxPool, "Max number of pool connections")
	cmd.Flags().String("babble.store_type", config.Babble.StoreType, "badger,inmem")
	cmd.Flags().String("babble.store_path", config.Babble.StorePath, "File containing the store database")
	viper.BindPFlags(cmd.Flags())
}

//NewBabbleCmd returns the command that starts EVM-Lite with Babble consensus
func NewBabbleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "babble",
		Short: "Run the evm-lite node with Babble consensus",
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {

			logger.WithFields(logrus.Fields{
				"Babble": config.Babble,
			}).Debug("Config")

			return nil
		},
		RunE: runBabble,
	}

	AddBabbleFlags(cmd)

	return cmd
}

func runBabble(cmd *cobra.Command, args []string) error {

	babble := babble.NewInmemBabble(*config.Babble, logger)
	engine, err := engine.NewEngine(*config, babble, logger)
	if err != nil {
		return fmt.Errorf("Error building Engine: %s", err)
	}

	engine.Run()

	return nil
}
