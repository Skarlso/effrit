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
	projectName   string
	parallelFiles int
}

func init() {
	scanCmd.Flags().StringVarP(&scanCmdOptions.projectName, "project-name", "p", "", "Define the name of the project.")
	scanCmd.Flags().IntVarP(&scanCmdOptions.parallelFiles, "parallel-files", "n", 5, "Define the maximum number of files processed in parallel. ")
	RootCmd.AddCommand(scanCmd)
}

func scan(cmd *cobra.Command, args []string) error {
	packages, err := pkg.Scan(scanCmdOptions.projectName, scanCmdOptions.parallelFiles)
	if err != nil {
		return err
	}
	packages.GatherDependedOnByCount()
	packages.CalculateInstability()
	packages.CalculateAbstractnessOfPackages()
	packages.CalculateDistance()
	packages.Display()
	packages.Dump()
	return nil
}
