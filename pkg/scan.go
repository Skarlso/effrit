package pkg

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"path/filepath"
)

type Package struct {
	Name string
	FullName string
	Imports []string
	ImportCount int
	DependedOnByCount int
}

func Scan() (map[string]*Package, error) {
	packages := make(map[string]*Package)
	// Format: [packageName] = {outSide import count}
	c := "go"
	args := []string{
		"list",
		"-f",
		"{{.ImportPath}} {{join .Imports \",\"}}",
		"./...",
	}
	cmd := exec.Command(c, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return packages, err
	}
	if err := cmd.Start(); err != nil {
		return packages, err
	}
	b, err := ioutil.ReadAll(stdout)
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
		p := &Package{
			Name: filepath.Base(string(pkg)),
			Imports: make([]string, 0),
			ImportCount: 0,
			DependedOnByCount: 0,
			FullName: string(pkg),
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
			}
		}
	}
	return packages, nil
}
