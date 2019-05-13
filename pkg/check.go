package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	EffritFileName = ".effrit_package_data.json"
)

func Check() error {
	if _, err := os.Stat(EffritFileName); err != nil {
		if os.IsNotExist(err) {
			return errors.New(EffritFileName)
		}
	}
	var packages = struct {
		Packages []Package `json:"packages"`
	}{
		Packages: make([]Package, 0),
	}
	data, err := ioutil.ReadFile(EffritFileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &packages)
	if err != nil {
		return err
	}
	fmt.Println(packages)
	return nil
}
