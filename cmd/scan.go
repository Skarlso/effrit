// @package_owner = @skarlso
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
	colorMode     bool
}

func init() {
	scanCmd.Flags().StringVarP(&scanCmdOptions.projectName, "scan-project-name", "p", "", "Define the name of the project.")
	scanCmd.Flags().IntVarP(&scanCmdOptions.parallelFiles, "scan-parallel-files", "n", 5, "Define the maximum number of files processed in parallel. ")
	scanCmd.Flags().BoolVarP(&scanCmdOptions.colorMode, "color", "c", false, "Color mode.")
	RootCmd.AddCommand(scanCmd)
}

func scan(cmd *cobra.Command, args []string) error {
	pkgs, err := pkg.Scan(scanCmdOptions.projectName, scanCmdOptions.parallelFiles)
	if err != nil {
		return err
	}
	pkgs.Display(scanCmdOptions.colorMode)
	pkgs.Dump()
	return nil
}
