package cmd

import "github.com/spf13/cobra"

// CheckCmd is the root for check commands.
var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check related commands, like check pr.",
}

func init() {
	RootCmd.AddCommand(CheckCmd)
}
