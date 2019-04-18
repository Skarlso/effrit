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

	"github.com/schollz/progressbar/v2"
)

func Scan(ignoreList []string) error {
	goFiles := make(map[string]string, 0)
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if contains(ignoreList, info.Name()) {
				return filepath.SkipDir
			}
		} else {
			if filepath.Ext(info.Name()) == ".go" {
				goFiles[path] = info.Name()
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	log.Println("Scanning total number of go files: ", len(goFiles))
	total := len(goFiles)
	bar := progressbar.New(total)
	//listOfPackages := make([]string, 0)
	for path, file := range goFiles {
		src, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, file, src, parser.PackageClauseOnly)
		if err != nil {
			return err
		}
		ast.Print(fset, f)
		ast.Inspect(f, func(n ast.Node) bool{
			var s string
			switch x := n.(type) {
			case *ast.Ident:
				s = x.Name
			}
			if s != "" {
				fmt.Println("Name: ", s)
				//fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
			}
			return true
		})
		//for _, i := range f.Imports {
		//	fmt.Println(i.Path.Value)
		//}
		_ = bar.Add(1)
	}
	_ = bar.Finish()
	return nil
}

func contains(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}