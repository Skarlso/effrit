package cmd

import (
	"github.com/Skarlso/effrit/pkg"
	"github.com/spf13/cobra"
	"strings"
)

var scanCmd = &cobra.Command{
	Use: "scan",
	Short: "scan a project for go files and packages",
	Long:`Scan will look in a folder for all go files and parse out the package
they are in. It will gather all defined packages.
`,
	RunE: scan,
}

var scanOptions struct {
	ignoreFolders string
}

func init() {
	scanCmd.Flags().StringVarP(&scanOptions.ignoreFolders,
		"ignore-folder",
		"i",
		"",
		"Comma separated list of folders to ignore.")
	RootCmd.AddCommand(scanCmd)
}

func scan(cmd *cobra.Command, args[] string) error {
	folders := make([]string, 0)
	if len(scanOptions.ignoreFolders) == 1 {
		folders = append(folders, scanOptions.ignoreFolders)
	} else if len(scanOptions.ignoreFolders) > 1 {
		folders = strings.Split(scanOptions.ignoreFolders, ",")
	}
	return pkg.Scan(folders)
}