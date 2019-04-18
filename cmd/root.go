package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use: "effrit",
	Short: "effrit -- look deep",
	SilenceUsage: true,
}
