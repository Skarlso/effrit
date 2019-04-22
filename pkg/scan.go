package pkg

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
)

// Package is a single package as determined by go list.
type Package struct {
	Name              string
	FullName          string
	Imports           []string
	ImportCount       int
	DependedOnByCount int
	Stability         float64
}

// Scan will scan a project using go list. As go list is running
// in the background, scan will display a waiting indicator.
func Scan() (map[string]Package, error) {
	packages := make(map[string]Package)
	// Format: [packageName] = {outSide import count}
	c := "go"
	args := []string{
		"list",
		"-f",
		"{{.ImportPath}} {{join .Imports \",\"}}",
		"./...",
	}
	cmd := exec.Command(c, args...)
	fmt.Println("Waiting for go list to finish scanning the project...")
	b, err := cmd.Output()
	if err != nil {
		return packages, err
	}
	lines := bytes.Split(b, []byte("\n"))
	for _, line := range lines {
		split := bytes.Split(line, []byte(" "))
		if len(split) < 2 {
			continue
		}

		pkg := split[0]
		imports := split[1]
		is := bytes.Split(imports, []byte(","))
		p := Package{
			Name:              filepath.Base(string(pkg)),
			Imports:           make([]string, 0),
			ImportCount:       0,
			DependedOnByCount: 0,
			FullName:          string(pkg),
			Stability:         0.0,
		}
		for _, i := range is {
			if bytes.Contains(i, []byte(".")) {
				p.Imports = append(p.Imports, string(i))
				p.ImportCount++
			}
		}
		packages[p.FullName] = p
	}
	for _, v := range packages {
		imports := v.Imports
		for _, i := range imports {
			if p, ok := packages[i]; ok {
				p.DependedOnByCount++
				packages[p.FullName] = p
			}
		}
	}
	return packages, nil
}
