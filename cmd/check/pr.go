package cmd

import (
	"github.com/Skarlso/effrit/cmd"
	"github.com/Skarlso/effrit/pkg"
	"github.com/spf13/cobra"
)

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Check a PR if it had changes added to a package.",
	Long: `Check will scan a project and compare changes to dependencies based on a locally stored file
called effrit_package_data.json. If there are changes to dependencies, Effrit will notify the person tagged
in a comment that new dependencies have been added or removed.`,
	RunE: check,
}

var prCmdOptions struct {
	projectName   string
	parallelFiles int
	repo          string
	prNumber      int
	owner         string
}

func init() {
	prCmd.Flags().StringVarP(&prCmdOptions.projectName, "project-name", "p", "", "Define the name of the project.")
	prCmd.Flags().IntVarP(&prCmdOptions.parallelFiles, "parallel-files", "n", 5, "Define the maximum number of files processed in parallel.")
	prCmd.Flags().StringVarP(&prCmdOptions.owner, "owner", "o", "", "The owner of the repository.")
	prCmd.Flags().StringVarP(&prCmdOptions.repo, "repo", "r", "", "The repository name to use to find the PR to comment on.")
	prCmd.Flags().IntVarP(&prCmdOptions.prNumber, "pr-number", "q", -1, "The PR to submit messages to in case that is needed.")
	cmd.CheckCmd.AddCommand(prCmd)
}

func check(cmd *cobra.Command, args []string) error {
	return pkg.Check(prCmdOptions.projectName, prCmdOptions.parallelFiles, prCmdOptions.owner, prCmdOptions.repo, prCmdOptions.prNumber)
}
