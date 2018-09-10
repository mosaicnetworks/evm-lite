package commands

import (
	"fmt"

	"github.com/mosaicnetworks/evm-lite/consensus/babble"
	"github.com/mosaicnetworks/evm-lite/engine"
	"github.com/spf13/cobra"
)

var babbleConfig = babble.DefaultConfig()

//AddBabbleFlags adds flags to the Babble command
func AddBabbleFlags(cmd *cobra.Command) {

	cmd.Flags().String("babble.proxy_addr", babbleConfig.ProxyAddr, "IP:PORT of Babble proxy")
	cmd.Flags().String("babble.client_addr", babbleConfig.ClientAddr, "IP:PORT to bind client proxy")
	cmd.Flags().String("babble.dir", babbleConfig.BabbleDir, "Directory contaning priv_key.pem and peers.json files")
	cmd.Flags().String("babble.node_addr", babbleConfig.NodeAddr, "IP:PORT of Babble node")
	cmd.Flags().String("babble.api_addr", babbleConfig.BabbleAPIAddr, "IP:PORT of Babble HTTP API service")
	cmd.Flags().Int("babble.heartbeat", babbleConfig.Heartbeat, "Heartbeat time milliseconds (time between gossips)")
	cmd.Flags().Int("babble.tcp_timeout", babbleConfig.TCPTimeout, "TCP timeout milliseconds")
	cmd.Flags().Int("babble.cache_size", babbleConfig.CacheSize, "Number of items in LRU caches")
	cmd.Flags().Int("babble.sync_limit", babbleConfig.SyncLimit, "Max number of Events per sync")
	cmd.Flags().Int("babble.max_pool", babbleConfig.MaxPool, "Max number of pool connections")
	cmd.Flags().String("babble.store_type", babbleConfig.StoreType, "badger,inmem")
	cmd.Flags().String("babble.store_path", babbleConfig.StorePath, "File containing the store database")
}

//NewBabbleCmd returns the command that starts EVM-Lite with Babble consensus
func NewBabbleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "babble",
		Short: "Run the evm-lite node with Babble consensus",
		RunE:  runBabble,
	}

	AddBabbleFlags(cmd)
	return cmd
}

func runBabble(cmd *cobra.Command, args []string) error {

	babble := babble.NewInmemBabble(*babbleConfig, logger)
	engine, err := engine.NewEngine(*config, babble, logger)
	if err != nil {
		return fmt.Errorf("Error building Engine: %s", err)
	}

	engine.Run()

	return nil
}
