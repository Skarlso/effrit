package pkg

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/schollz/progressbar/v2"
)

func Scan(ignoreList []string) error {
	goFiles := make([]os.FileInfo, 0)
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
				goFiles = append(goFiles, info)
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
	for _ = range goFiles {
		_ = bar.Add(1)
		time.Sleep(time.Millisecond)
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