package main

import (
	cmd "github.com/mosaicnetworks/evm-lite/cmd/evml/commands"
)

func main() {

	rootCmd := cmd.RootCmd

	rootCmd.AddCommand(
		cmd.NewSoloCmd(),
		cmd.NewBabbleCmd(),
		cmd.VersionCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
