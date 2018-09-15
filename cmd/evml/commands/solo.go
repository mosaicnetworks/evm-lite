package commands

import (
	"fmt"

	"github.com/mosaicnetworks/evm-lite/consensus/solo"
	"github.com/mosaicnetworks/evm-lite/engine"
	"github.com/spf13/cobra"
)

//NewSoloCmd returns the command that starts EVM-Lite with Solo consensus
func NewSoloCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "solo",
		Short: "Run the evm-lite node with Solo consensus (no consensus)",
		RunE:  runSolo,
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
