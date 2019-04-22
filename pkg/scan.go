package pkg

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
)

func Scan() error {
	// Format: [packageName] = {outSide import count}
	packages := make(map[string]int)
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
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	b, err := ioutil.ReadAll(stdout)
	if err != nil {
		return err
	}
	lines := bytes.Split(b, []byte("\n"))
	for _, line := range lines {
		split := bytes.Split(line, []byte(" "))
		if len(split) < 2 {
			continue
		}
		pkg := filepath.Base(string(split[0]))
		imports := split[1]
		is := bytes.Split(imports, []byte(","))
		packages[pkg] = 0
		for _, i := range is {
			if bytes.Contains(i, []byte(".")) {
				packages[pkg]++
			}
		}
	}
	fmt.Println(packages)
	return nil
}
