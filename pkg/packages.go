package pkg

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
)

// Package is a single package as determined by go list.
type Package struct {
	Name               string
	FullName           string
	Imports            []string
	ImportCount        float64
	DependedOnByCount  float64
	Stability          float64
	Abstractness       float64
	DistanceFromMedian float64
	Dir                string
	GoFiles            []string
}

// Packages represents a collection of packages.
type Packages struct {
	packageMap map[string]Package
}

// NewPackages returns a new initialized package collection.
func NewPackages() *Packages {
	return &Packages{
		packageMap: make(map[string]Package),
	}
}

func (pkg *Packages) gatherDependedOnByCount() {
	for _, v := range pkg.packageMap {
		imports := v.Imports
		for _, i := range imports {
			if p, ok := pkg.packageMap[i]; ok {
				p.DependedOnByCount++
				pkg.packageMap[p.FullName] = p
			}
		}
	}
}

// ParseGoFiles will walk all the files in the package
// and analyies abstractness.
func (pkg *Packages) ParseGoFiles() {
	for _, p := range pkg.packageMap {
		for _, f := range p.GoFiles {
			fset := token.NewFileSet()
			data, err := ioutil.ReadFile(filepath.Join(p.Dir, f))
			if err != nil {
				log.Fatal(err)
			}
			node, err := parser.ParseFile(fset, f, data, 0)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

var yellow = color.New(color.FgYellow)
var red = color.New(color.FgRed)
var green = color.New(color.FgGreen, color.Bold)

// Display displays the analysed information in a pretty way...
// TODO: Add multiple display options and Graph generation.
func (pkg *Packages) Display() {
	const padding = 3
	table := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	writeColumns(table, []string{"NAME", "STABILITY"})
	for _, p := range pkg.packageMap {
		c := &color.Color{}
		if p.Stability < 0.5 {
			c = green
		} else if p.Stability >= 0.5 && p.Stability < 1 {
			c = yellow
		} else if p.Stability == 1 {
			c = red
		}
		stability := fmt.Sprintf("%.1f", p.Stability)
		writeColumns(table, []string{p.FullName, c.Sprint(stability)})
	}
	_ = table.Flush()
}

func writeColumns(w io.Writer, cols []string) {
	_, _ = fmt.Fprintln(w, strings.Join(cols, "\t"))
}
