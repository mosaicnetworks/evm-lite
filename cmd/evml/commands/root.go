package commands

import (
	"github.com/mosaicnetworks/evm-lite/cmd/evml/commands/keys"
	"github.com/mosaicnetworks/evm-lite/cmd/evml/commands/run"
	"github.com/spf13/cobra"
)

//RootCmd is the root command for evml
var RootCmd = &cobra.Command{
	Use:   "evml",
	Short: "EVM-Lite",
}

func init() {
	RootCmd.AddCommand(
		run.RunCmd,
		keys.KeysCmd,
		VersionCmd,
	)
	//do not print usage when error occurs
	RootCmd.SilenceUsage = true
}
