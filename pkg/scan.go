package pkg

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Scan will scan a project using go list. As go list is running
// in the background, scan will display a waiting indicator.
func Scan(projectName string, parallel int) (*Packages, error) {
	pkgs := NewPackages(parallel)
	// Format: [packageName] = {outSide import count}
	c := "go"
	args := []string{
		"list",
		"-f",
		"{{.ImportPath}} {{join .Imports \",\"}} {{join .GoFiles \",\"}} {{.Dir}}",
		"./...",
	}
	/* #nosec */
	cmd := exec.Command(c, args...)
	fmt.Print("Waiting for go list to finish scanning the project...")
	b, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	fmt.Print("done.\n")
	lines := bytes.Split(b, []byte("\n"))
	for _, line := range lines {
		split := bytes.Split(line, []byte(" "))
		if len(split) < 2 {
			continue
		}

		pkg := split[0]
		imports := split[1]
		goFiles := string(split[2])
		dir := split[3]
		gf := strings.Split(goFiles, ",")

		is := bytes.Split(imports, []byte(","))
		p := Package{
			Name:              filepath.Base(string(pkg)),
			Imports:           make([]string, 0),
			ImportCount:       0,
			DependedOnByCount: 0,
			FullName:          string(pkg),
			Stability:         0.0,
			GoFiles:           gf,
			Dir:               string(dir),
		}
		for _, i := range is {
			if bytes.Contains(i, []byte(".")) {
				if len(projectName) > 0 && !bytes.Contains(i, []byte(projectName)) {
					continue
				}
				p.Imports = append(p.Imports, string(i))
				p.ImportCount++
			}
		}
		pkgs.packageMap[p.FullName] = p
		pkgs.packageNames = append(pkgs.packageNames, p.FullName)
	}
	sort.Strings(pkgs.packageNames)
	pkgs.GatherDependedOnByCount()
	pkgs.CalculateInstability()
	err = pkgs.CalculateAbstractnessOfPackages()
	if err != nil {
		return pkgs, err
	}
	pkgs.CalculateDistance()
	return pkgs, nil
}
