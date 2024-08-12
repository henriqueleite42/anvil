package cmd

import (
	"fmt"

	"github.com/anvil/anvil/cmd/config"
	"github.com/spf13/cobra"
)

func addVersionCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number of anvil",
		Long:  `All software has versions. This is anvil's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("anvil Micro-Services Generator " + config.CLI_VERSION + " -- HEAD")
		},
	})
}
