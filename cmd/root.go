package cmd

import "github.com/spf13/cobra"

// RootCmd is the root cobra command.
var RootCmd = &cobra.Command{
	Use:          "effrit",
	Short:        "effrit -- look deep",
	SilenceUsage: true,
}
