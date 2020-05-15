// @package_owner = @skarlso
package cmd

import (
	"os"

	"github.com/Skarlso/effrit/pkg"
	"github.com/Skarlso/effrit/pkg/providers/github"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Create a server which listens for PRs and does a check-pr on them.",
	Long:  `Effrit can be run as a server which can be registered as a webhook for PR events.`,
	RunE:  run,
}

var serverCmdOptions struct {
	pkg.ServerConfig
}

func init() {
	serverCmd.Flags().BoolVar(&serverCmdOptions.AutoTLS, "auto-tls", false, "--auto-tls")
	serverCmd.Flags().StringVar(&serverCmdOptions.CacheDir, "cache-dir", "", "--cache-dir /home/user/.server/.cache")
	serverCmd.Flags().StringVar(&serverCmdOptions.ServerKeyPath, "server-key-path", "", "--server-key-file /home/user/.server/server.key")
	serverCmd.Flags().StringVar(&serverCmdOptions.ServerCrtPath, "server-crt-path", "", "--server-crt-file /home/user/.server/server.crt")
	serverCmd.Flags().StringVar(&serverCmdOptions.Port, "port", "9998", "--port 443")
	serverCmd.Flags().StringVar(&serverCmdOptions.Hostname, "hostname", "0.0.0.0", "--hostname pr-check.example.com")
	serverCmd.Flags().StringVarP(&serverCmdOptions.ProjectName, "project-name", "p", "", "Define the name of the project.")
	serverCmd.Flags().IntVarP(&serverCmdOptions.Parallel, "parallel-files", "n", 5, "Define the maximum number of files processed in parallel.")
	RootCmd.AddCommand(serverCmd)
}

func run(cmd *cobra.Command, args []string) error {
	// Create the commenter
	githubCommenter := github.NewGithubCommenter(github.Config{
		ProjectName: serverCmdOptions.ProjectName,
		Parallel:    serverCmdOptions.Parallel,
	}, github.Dependencies{})
	// Create the server
	server := pkg.NewServer(serverCmdOptions.ServerConfig, pkg.Dependencies{
		Commenter: githubCommenter,
		Logger:    zerolog.New(os.Stderr),
	})

	return server.Run()
}
