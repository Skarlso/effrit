package pkg

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

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
	DependedOnByNames  []string
	Stability          float64
	Abstractness       float64
	DistanceFromMedian float64
	Dir                string
	GoFiles            []string
}

// Packages represents a collection of packages.
type Packages struct {
	packageMap   map[string]Package
	packageNames []string
	semaphore    int
}

// NewPackages returns a new initialized package collection.
func NewPackages(parallel int) *Packages {
	return &Packages{
		packageMap:   make(map[string]Package),
		packageNames: make([]string, 0),
		semaphore:    parallel,
	}
}

// GatherDependedOnByCount looks for packages which import other packages in the project.
func (pkg *Packages) GatherDependedOnByCount() {
	for _, v := range pkg.packageMap {
		imports := v.Imports
		for _, i := range imports {
			if p, ok := pkg.packageMap[i]; ok {
				p.DependedOnByCount++
				if len(p.DependedOnByNames) < 1 {
					p.DependedOnByNames = make([]string, 0)
				}
				p.DependedOnByNames = append(p.DependedOnByNames, v.FullName)
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
		var wg sync.WaitGroup
		funcCount := 0.0
		abstractsCount := 0.0
		fmt.Printf("Scanning %s go file(s) in package %s.\n", keyName.Sprint(len(p.GoFiles)), keyName.Sprint(p.FullName))
		sem := make(chan int, pkg.semaphore)
		errChan := make(chan error, len(p.GoFiles))
		funcChan := make(chan float64, len(p.GoFiles))
		absChan := make(chan float64, len(p.GoFiles))
		for _, f := range p.GoFiles {
			wg.Add(1)
			go parseGoFile(f, p.Dir, funcChan, absChan, sem, errChan, &wg)
		} // go files in packages
		wg.Wait()
		// Need to close the channels here to be able to for on them.
		close(funcChan)
		close(absChan)
		close(errChan)

		errorList := make([]error, 0)
		for e := range errChan {
			errorList = append(errorList, e)
		}
		if len(errorList) > 0 {
			fmt.Printf("%d error(s) processing pkg %s\n", len(errorList), p.FullName)
			fmt.Println("listing error(s):")
			for _, e := range errorList {
				fmt.Println(e)
			}
			fmt.Println("Please fix these before continuing.")
			os.Exit(1)
		}
		for a := range absChan {
			abstractsCount += a
		}
		for a := range funcChan {
			funcCount += a
		}
		p.Abstractness = abstractsCount / funcCount
		pkg.packageMap[k] = p
	} // packages
}

func parseGoFile(fh, dir string,
	funcChan, absChan chan float64,
	sem chan int,
	errChan chan error,
	wg *sync.WaitGroup) {

	defer wg.Done()
	sem <- 1
	defer func() { <-sem }()
	funcCount := 0.0
	abstractsCount := 0.0
	fset := token.NewFileSet()
	if len(fh) < 1 {
		fmt.Println("skipping folder... file is empty:", dir)
		return
	}
	/* #nosec */
	data, err := ioutil.ReadFile(filepath.Join(dir, fh))
	if err != nil {
		errChan <- err
		return
	}
	node, err := parser.ParseFile(fset, fh, data, 0)
	if err != nil {
		errChan <- err
		return
	}
	ast.Inspect(node, func(n ast.Node) bool {
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
	funcChan <- funcCount
	absChan <- abstractsCount
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
	for _, pname := range pkg.packageNames {
		p := pkg.packageMap[pname]
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

// Dump dumps the connections between packages in a JSON format
// for other tools to process.
func (pkg *Packages) Dump() {
	pkgs := make([]Package, 0)
	for _, p := range pkg.packageMap {
		pkgs = append(pkgs, p)
	}
	var packages = struct {
		Packages []Package `json:"packages"`
	}{
		Packages: pkgs,
	}
	data, err := json.Marshal(packages)
	if err != nil {
		log.Fatal("error marshaling to json: ", err)
	}
	err = ioutil.WriteFile(".effrit_package_data.json", data, 0666)
	if err != nil {
		log.Fatal("error while writing json file: ", err)
	}
}
