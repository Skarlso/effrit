package cmd

import (
	"github.com/Skarlso/effrit/pkg"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check-pr",
	Short: "Check a PR if it had changes added to a package.",
	Long: `Check will scan a project and compare changes to dependencies based on a locally stored file
called effrit_package_data.json. If there are changes to dependencies, Effrit will notify the person tagged
in a comment that new dependencies have been added or removed.`,
	RunE: check,
}

var checkCmdOptions struct {
	projectName   string
	parallelFiles int
}

func init() {
	checkCmd.Flags().StringVarP(&checkCmdOptions.projectName, "project-name", "p", "", "Define the name of the project.")
	checkCmd.Flags().IntVarP(&checkCmdOptions.parallelFiles, "parallel-files", "n", 5, "Define the maximum number of files processed in parallel. ")
	RootCmd.AddCommand(checkCmd)
}

func check(cmd *cobra.Command, args []string) error {
	return pkg.Check()
}
