// Package pkg has the main brains of Effrit.
// @package_owner = @skarlso
package pkg

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
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
	data, err := ioutil.ReadFile(EffritFileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &packages)
	if err != nil {
		return err
	}

	// Construct a map of package names and their imports.
	packageMap := make(map[string][]string)

	for _, p := range packages.Packages {
		sort.Strings(p.DependedOnByNames)
		packageMap[p.FullName] = p.DependedOnByNames
	}

	_, err = Scan(projectName, parallel)
	if err != nil {
		return err
	}
	packages.Packages = make([]Package, 0)
	data, _ = ioutil.ReadFile(EffritFileName)
	_ = json.Unmarshal(data, &packages)

	ownersToContact := make([]string, 0)
	// Compare the new result with the old result's map data.
	for _, p := range packages.Packages {
		dependents := packageMap[p.FullName]
		sort.Strings(dependents)
		sort.Strings(p.DependedOnByNames)
		if !reflect.DeepEqual(dependents, p.DependedOnByNames) {
			fmt.Println("A new dependency has been added to package: ", p.FullName)
			owner, err := getOwnerForFile(p.Dir, p.GoFiles)
			if err != nil {
				return err
			}
			ownersToContact = append(ownersToContact, owner)
		}
	}
	if len(ownersToContact) > 0 {
		return contactOwners(ownersToContact, owner, repo, prNumber)
	}
	fmt.Println("Everything is the same. Carry on.")
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

func contactOwners(owners []string, owner, repo string, n int) error {
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
	com := "This is my comment"
	comment := github.PullRequestComment{
		Body: &com,
	}
	_, _, err := client.PullRequests.CreateComment(ctx, owner, repo, n, &comment)
	return err
}
