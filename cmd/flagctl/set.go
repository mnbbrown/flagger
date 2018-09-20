package main

import (
	"fmt"
	"os"

	flagger "github.com/mnbbrown/flagger/client"
	"github.com/spf13/cobra"
)

var setCommand = &cobra.Command{
	Use: "set",
	Run: func(cmd *cobra.Command, args []string) {
		client := flagger.NewClient("http://localhost:8082")
		switch len(args) {
		case 3:
			fmt.Printf("Setting flag %s", args[0])
			client.Set(args[0], "default", args[2], args[3])
		case 4:
			fmt.Printf("Setting flag %s (environment: %s)", args[0], args[1])
			client.Set(args[0], args[1], args[2], args[3])
		default:
			fmt.Println("Incorrect number of arguments")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(setCommand)
}
