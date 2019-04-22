package pkg

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

// Package is a single package as determined by go list.
type Package struct {
	Name              string
	FullName          string
	Imports           []string
	ImportCount       int
	DependedOnByCount int
	// For stability, 0.0 is a vaild value. Hence we need a value
	// where stability has not yet been calculated for a given
	// package.
	Stability *float64
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
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return packages, err
	}
	if err := cmd.Start(); err != nil {
		return packages, err
	}
	waitForCommand(cmd)

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
		p := Package{
			Name:              filepath.Base(string(pkg)),
			Imports:           make([]string, 0),
			ImportCount:       0,
			DependedOnByCount: 0,
			FullName:          string(pkg),
			Stability:         nil,
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

func waitForCommand(cmd *exec.Cmd) {
	var wg sync.WaitGroup
	wg.Add(1)
	spinner := `|/-\`
	done := make(chan struct{}, 0)
	// TODO: think of a reasonable timeout value. Context with Timeout

	go func() {
		defer wg.Done()
		cmd.Wait()
		done <- struct{}{}
	}()
	go func() {
		counter := 0
		for {
			select {
			case <-done:
				break
			default:
			}
			counter = (counter + 1) % len(spinner)
			fmt.Printf("\r[%s] Waiting for command to finish...", string(spinner[counter]))
			time.Sleep(500 * time.Millisecond)
		}
	}()
	wg.Wait()
}
