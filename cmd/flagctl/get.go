package main

import (
	"fmt"
	"os"

	flagger "github.com/mnbbrown/flagger/client"
	"github.com/spf13/cobra"
)

var getCommand = &cobra.Command{
	Use: "get",
	Run: func(cmd *cobra.Command, args []string) {
		client := flagger.NewClient("http://localhost:8082")
		switch len(args) {
		case 1:
			value := client.Get(args[0], "default")
			fmt.Println(value)
		case 2:
			value := client.Get(args[0], args[1])
			fmt.Println(value)
		default:
			fmt.Println("Incorrect number of arguments")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCommand)
}
