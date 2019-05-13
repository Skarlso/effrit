// Some comment
// @package_author = @skarlso
package pkg

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

const (
	EffritFileName = ".effrit_package_data.json"
	CommentTag     = "@package_author"
)

func Check(projectName string, parallel int) error {
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

	err = Scan(projectName, parallel)
	if err != nil {
		return err
	}
	packages.Packages = make([]Package, 0)
	data, _ = ioutil.ReadFile(EffritFileName)
	_ = json.Unmarshal(data, &packages)

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
			fmt.Println("Contacting owner: ", owner)
		}
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
				_, _ = fmt.Sscanf(line, "// @package_author = %s", &owner)
				_ = file.Close()
				return owner, nil
			}
		}
		_ = file.Close()
	}
	return "", nil
}
