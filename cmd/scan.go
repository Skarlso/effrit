package cmd

import (
	"github.com/Skarlso/effrit/pkg"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scan a project for go files and packages",
	Long: `Scan will look in a folder for all go files and parse out the package
they are in. It will gather all defined packages.
`,
	RunE: scan,
}

var scanCmdOptions struct {
	projectName string
}

func init() {
	scanCmd.Flags().StringVarP(&scanCmdOptions.projectName, "project-name", "p", "", "Define the name of the project.")
	RootCmd.AddCommand(scanCmd)
}

func scan(cmd *cobra.Command, args []string) error {
	packages, err := pkg.Scan(scanCmdOptions.projectName)
	if err != nil {
		return err
	}
	packages = pkg.Analyse(packages)
	pkg.Display(packages)
	return nil
}
