// Some comment
// @package_author = Gergely Brautigam
package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
)

const (
	EffritFileName = ".effrit_package_data.json"
	CommentSection = "@package_author = %s"
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

	data, _ = ioutil.ReadFile(EffritFileName)
	_ = json.Unmarshal(data, &packages)

	// Compare the new result with the old result's map data.
	for _, p := range packages.Packages {
		dependents := packageMap[p.Name]
		sort.Strings(dependents)
		sort.Strings(p.DependedOnByNames)
		if !reflect.DeepEqual(dependents, p.DependedOnByNames) {
			fmt.Println("A new dependency has been added to package: ", p.Name)
		}
	}
	return nil
}
