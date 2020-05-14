// @package_owner = @skarlso
package main

import (
	"log"

	"github.com/Skarlso/effrit/cmd"
	_ "github.com/Skarlso/effrit/cmd/check"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
