package pkg

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGatherDependentOnByCount(t *testing.T) {
	pkgs := NewPackages(1)
	testMap := map[string]Package{
		"github.com/Test/Test": {
			Name:               "Test",
			FullName:           "github.com/Test/Test",
			ImportCount:        0.0,
			Stability:          0.0,
			Abstractness:       0.0,
			DistanceFromMedian: 0.0,
			Dir:                "/home/go/test",
			GoFiles:            []string{"main.go", "other.go"},
		},
		"github.com/Test/Test2": {
			Name:               "Test2",
			FullName:           "github.com/Test/Test2",
			Imports:            []string{"github.com/Test/Test"},
			ImportCount:        1.0,
			Stability:          0.0,
			Abstractness:       0.0,
			DistanceFromMedian: 0.0,
			Dir:                "/home/go/test",
			GoFiles:            []string{"other2.go"},
		},
	}
	pkgs.packageMap = testMap
	pkgs.GatherDependedOnByCount()
	if pkgs.packageMap["github.com/Test/Test"].DependedOnByCount < 1 {
		t.Fatal("dependedOnByCount should have been greater than 0")
	}
	pkgs.CalculateInstability()
	if pkgs.packageMap["github.com/Test/Test"].Stability != 0.0 {
		t.Fatal("Stability of package test should have been 1.0. was: ", pkgs.packageMap["github.com/Test/Test"].Stability)
	}
	if pkgs.packageMap["github.com/Test/Test2"].Stability != 1.0 {
		t.Fatal("Stability of package test2 should have been 0.0. was: ", pkgs.packageMap["github.com/Test/Test2"].Stability)
	}
}

func TestCalculateAbstractness(t *testing.T) {
	pkgs := NewPackages(1)
	tmp, _ := ioutil.TempDir("", "TestCalculateAbstractness")
	main := `package test
type s interface {
	FunkyTime() error
}
`
	other := `package test2
func someFunction() {

}
`
	_ = os.Mkdir(filepath.Join(tmp, "test"), 0777)
	_ = os.Mkdir(filepath.Join(tmp, "test2"), 0777)
	_ = ioutil.WriteFile(filepath.Join(tmp, "test", "main.go"), []byte(main), 0777)
	_ = ioutil.WriteFile(filepath.Join(tmp, "test2", "other.go"), []byte(other), 0777)
	testMap := map[string]Package{
		"github.com/Test/Test": {
			Name:               "Test",
			FullName:           "github.com/Test/Test",
			Imports:            []string{"github.com/Test/Test3"},
			ImportCount:        1.0,
			Stability:          0.0,
			Abstractness:       0.0,
			DistanceFromMedian: 0.0,
			Dir:                filepath.Join(tmp, "test"),
			GoFiles:            []string{"main.go"},
		},
		"github.com/Test/Test2": {
			Name:               "Test2",
			FullName:           "github.com/Test/Test2",
			Imports:            []string{"github.com/Test/Test"},
			ImportCount:        1.0,
			Stability:          0.0,
			Abstractness:       0.0,
			DistanceFromMedian: 0.0,
			Dir:                filepath.Join(tmp, "test2"),
			GoFiles:            []string{"other.go"},
		},
	}
	pkgs.packageMap = testMap
	pkgs.GatherDependedOnByCount()
	err := pkgs.CalculateAbstractnessOfPackages()
	if err != nil {
		t.Fatal(err)
	}
	if pkgs.packageMap["github.com/Test/Test"].Abstractness != 0.5 {
		t.Fatal("Abstractness of package test should have been 0.5. was: ", pkgs.packageMap["github.com/Test/Test"].Abstractness)
	}
	if pkgs.packageMap["github.com/Test/Test2"].Abstractness != 0.0 {
		t.Fatal("Abstractness of package test2 should have been 0.0. was: ", pkgs.packageMap["github.com/Test/Test2"].Abstractness)
	}
}
