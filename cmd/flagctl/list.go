package main

import (
	"fmt"
	"os"

	flagger "github.com/mnbbrown/flagger/client"
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		client := flagger.NewClient("http://localhost:8082")
		flags, err := client.List()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for flag, envs := range flags {
			for env, value := range envs {
				if env == "default" {
					fmt.Printf("%s: %s %v\n", flag, value.Type, value.InternalValue)
				} else {
					fmt.Printf("%s (env %s): %s %v\n", flag, env, value.Type, value.InternalValue)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCommand)
}
