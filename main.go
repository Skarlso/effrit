package main

import (
	"github.com/spf13/cobra"
	"log"
)

var RootCmd = &cobra.Command{
	Use: "effrit",
	Short: "effrit -- look deep",
	SilenceUsage: true,
}

func main()  {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
