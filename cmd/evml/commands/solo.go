package commands

import (
	"fmt"

	"github.com/mosaicnetworks/evm-lite/src/consensus/solo"
	"github.com/mosaicnetworks/evm-lite/src/engine"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//NewSoloCmd returns the command that starts EVM-Lite with Solo consensus
func NewSoloCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "solo",
		Short: "Run the evm-lite node with Solo consensus (no consensus)",
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {

			config.SetDataDir(config.BaseConfig.DataDir)

			logger.WithFields(logrus.Fields{
				"Eth": config.Eth,
			}).Debug("Config")

			return nil
		},
		RunE: runSolo,
	}
	return cmd
}

func runSolo(cmd *cobra.Command, args []string) error {

	solo := solo.NewSolo(logger)
	engine, err := engine.NewEngine(*config, solo, logger)
	if err != nil {
		return fmt.Errorf("Error building Engine: %s", err)
	}

	engine.Run()

	return nil
}
