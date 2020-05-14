package github

import (
	"github.com/Skarlso/effrit/pkg"
	"github.com/Skarlso/effrit/pkg/providers"
)

// Config provides configuration for the github commenter.
type Config struct {
	ProjectName string
	Parallel    int
}

// Dependencies provides dependencies for the github commenter.
type Dependencies struct{}

type githubCommenter struct {
	Config
	Dependencies
}

// NewGithubCommenter returns a new github based pr commenter.
func NewGithubCommenter(cfg Config, deps Dependencies) providers.Commenter {
	return &githubCommenter{
		Config:       cfg,
		Dependencies: deps,
	}
}

// Comment will create a comment on a PR.
func (p *githubCommenter) Comment(owner string, repo string, number int) error {
	pkg.Check(p.Config.ProjectName, p.Config.Parallel, owner, repo, number)
	return nil
}
