package main

import (
	cmd "github.com/mosaicnetworks/evm-lite/cmd/evml/commands"
)

func main() {

	rootCmd := cmd.RootCmd

	rootCmd.AddCommand(
		cmd.NewSoloCmd(),
		cmd.NewBabbleCmd(),
		cmd.NewRaftCmd(),
		cmd.VersionCmd)

	//Do not print usage when error occurs
	rootCmd.SilenceUsage = true

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
