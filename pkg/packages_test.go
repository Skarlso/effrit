package pkg

import (
	"testing"
)

func TestGatherDependentOnByCount(t *testing.T) {
	pkgs := NewPackages(1)
	testMap := map[string]Package{
		"github.com/Test/Test": {
			Name:               "Test",
			FullName:           "github.com/Test/Test",
			Imports:            []string{"github.com/Test/Test3"},
			ImportCount:        1,
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
			ImportCount:        1,
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
}
