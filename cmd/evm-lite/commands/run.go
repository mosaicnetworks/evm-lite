package commands

import (
	"fmt"

	"github.com/mosaicnetworks/evm-lite/consensus/solo"
	"github.com/mosaicnetworks/evm-lite/engine"
	"github.com/spf13/cobra"
)

// NewRunCmd returns the command that allows the CLI to start a node.
func NewRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the evm-lite node",
		RunE:  run,
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) error {

	solo := solo.NewSolo(logger)
	engine, err := engine.NewEngine(*config, solo, logger)
	if err != nil {
		return fmt.Errorf("Error building Engine: %s", err)
	}

	engine.Run()

	return nil
}
