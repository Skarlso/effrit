package main

import (
	"github.com/Skarlso/effrit/cmd"
	"log"
)

func main()  {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
