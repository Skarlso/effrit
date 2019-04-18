package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use: "scan",
	Short: "scan a project for go files and packages",
	Long:`Scan will look in a folder for all go files and parse out the package
they are in. It will gather all defined packages.
`,
	RunE: scan,
}

func init() {
	RootCmd.AddCommand(scanCmd)
}

func scan(cmd *cobra.Command, args[] string) error {
	fmt.Println("yisss")
	return nil
}