// Package pkg has the main brains of Effrit.
// @package_owner = @skarlso
package pkg

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	// EffritFileName sets the name of the file Effrit is using to save package data.
	EffritFileName = ".effrit_package_data.json"
	// CommentTag is the name of the tag.
	CommentTag = "@package_owner"
	// CommentFormat is the format used to parse a line for owner name.
	CommentFormat = "// " + CommentTag + " = %s"
)

// Check will run an analysis of packages and detect if a new dependency has been added
// to package.
func Check(projectName string, parallel int, owner, repo string, prNumber int) error {
	if _, err := os.Stat(EffritFileName); err != nil {
		if os.IsNotExist(err) {
			return errors.New(EffritFileName)
		}
	}
	var packages = struct {
		Packages []Package `json:"packages"`
	}{
		Packages: make([]Package, 0),
	}
	data, err := os.ReadFile(EffritFileName)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &packages); err != nil {
		return err
	}

	// Construct a map of package names and their imports.
	packageMap := make(map[string][]string)

	for _, p := range packages.Packages {
		sort.Strings(p.DependedOnByNames)
		packageMap[p.FullName] = p.DependedOnByNames
	}

	pkgs, err := Scan(projectName, parallel)
	if err != nil {
		return err
	}
	pkgs.Dump()
	packages.Packages = make([]Package, 0)
	data, _ = os.ReadFile(EffritFileName)
	_ = json.Unmarshal(data, &packages)

	ownersToContact := make(map[string]string)
	// Compare the new result with the old result's map data.
	for _, p := range packages.Packages {
		dependents := packageMap[p.FullName]
		for _, dep := range dependents {
			if !contains(p.DependedOnByNames, dep) {
				owner, err := getOwnerForFile(p.Dir, p.GoFiles)
				if err != nil {
					return err
				}
				ownersToContact[owner] = dep
				break
			}
		}
	}
	if len(ownersToContact) > 0 {
		fmt.Print("Contacting owners about package dependency changes...")
		if err := contactOwners(ownersToContact, owner, repo, prNumber); err != nil {
			return err
		}
		fmt.Println("done.")
	}
	return nil
}

func getOwnerForFile(dir string, files []string) (string, error) {
	// we check until we find an owner tag for this package.
	for _, f := range files {
		file, err := os.Open(filepath.Join(dir, f))
		if err != nil {
			return "", err
		}
		fs := bufio.NewScanner(file)
		for fs.Scan() {
			line := fs.Text()
			if strings.Contains(line, CommentTag) {
				var owner string
				_, _ = fmt.Sscanf(line, CommentFormat, &owner)
				_ = file.Close()
				return owner, nil
			}
		}
		_ = file.Close()
	}
	return "", nil
}

func contactOwners(notify map[string]string, owner, repo string, n int) error {
	token := os.Getenv("EFFRIT_GITHUB_TOKEN")
	if len(token) < 1 {
		return errors.New("please define EFFRIT_GITHUB_TOKEN in order to allow effrit to create comments")
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	pr, _, _ := client.PullRequests.Get(ctx, owner, repo, n)
	c := "# Dependency Changed\n"
	for o, dep := range notify {
		c += fmt.Sprintf("%s please review your package. **%s** has been added as a new dependency\n", o, dep)
	}
	comment := github.IssueComment{
		AuthorAssociation: pr.AuthorAssociation,
		Body:              &c,
	}
	if _, _, err := client.Issues.CreateComment(ctx, owner, repo, n, &comment); err != nil {
		return err
	}
	return nil
}

func contains(a []string, elem string) bool {
	for _, e := range a {
		if e == elem {
			return true
		}
	}
	return false
}
