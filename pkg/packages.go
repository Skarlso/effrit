package pkg

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
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

// GatherDependedOnByCount looks for packages which import other packages in the project.
func (pkg *Packages) GatherDependedOnByCount() {
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

// CalculateInstability will calculate the instability metric for all packages.
func (pkg *Packages) CalculateInstability() {
	for k, v := range pkg.packageMap {
		v.Stability = v.ImportCount / (v.ImportCount + v.DependedOnByCount)
		pkg.packageMap[k] = v
	}
}

// CalculateDistance will calculate the distance from the main sequence for all packages.
// 0: As far away from the main sequence as possible
// 1: Close as possible
// 0,0: Zone of Pain
// 1,1: Zone of Uselessness
func (pkg *Packages) CalculateDistance() {
	for k, v := range pkg.packageMap {
		v.DistanceFromMedian = abs(v.Stability + v.Abstractness - 1)
		pkg.packageMap[k] = v
	}
}

func abs(a float64) float64 {
	if a < 0 {
		return -a
	}
	return a
}

// CalculateAbstractnessOfPackages will walk all the files in the package
// and analyses abstractness.
func (pkg *Packages) CalculateAbstractnessOfPackages() {
	for k, p := range pkg.packageMap {
		funcCount := 0.0
		abstractsCount := 0.0
		fmt.Printf("Scanning %s go file(s) in package %s.\n", keyName.Sprint(len(p.GoFiles)), keyName.Sprint(p.FullName))
		for _, f := range p.GoFiles {
			fset := token.NewFileSet()
			/* #nosec */
			data, err := ioutil.ReadFile(filepath.Join(p.Dir, f))
			if err != nil {
				log.Fatal(err)
			}
			node, err := parser.ParseFile(fset, f, data, 0)
			if err != nil {
				log.Fatal(err)
			}
			ast.Inspect(node, func(n ast.Node) bool {
				// TODO: This needs structs with no receivers counted as interface.
				switch n.(type) {
				case *ast.FuncDecl:
					funcCount++
				case *ast.InterfaceType:
					abstractsCount++
				case *ast.StructType:
					// This right now calculates structs towards abstractness.
					// I have no easy way to find receivers for structs yet
					// so I'm counting all structs towards interfaces. If there are
					// implementations in these packages they would even out this count.
					abstractsCount++
				}
				return true
			})
		} // go files in packages
		p.Abstractness = abstractsCount / funcCount
		pkg.packageMap[k] = p
	} // packages
}

var keyName = color.New(color.FgWhite, color.Bold)
var yellow = color.New(color.FgYellow)
var red = color.New(color.FgRed)
var green = color.New(color.FgGreen, color.Bold)

// Display displays the analysed information in a pretty way...
// TODO: Add multiple display options and Graph generation.
func (pkg *Packages) Display() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAME", "STABILITY", "ABSTRACTNESS", "DISTANCE"})
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
		abstractness := fmt.Sprintf("%.1f", p.Abstractness)
		distance := fmt.Sprintf("%.1f", p.DistanceFromMedian)
		table.Append([]string{p.FullName, c.Sprint(stability), keyName.Sprint(abstractness), keyName.Sprint(distance)})
	}
	table.Render()
}
