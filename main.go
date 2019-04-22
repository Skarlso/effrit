package main

import (
	"log"

	"github.com/Skarlso/effrit/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
