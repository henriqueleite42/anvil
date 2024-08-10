package cmd

import (
	"fmt"

	"github.com/anvlet/anvlet/cmd/config"
	"github.com/spf13/cobra"
)

func addVersionCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number of anvlet",
		Long:  `All software has versions. This is anvlet's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("anvlet Micro-Services Generator " + config.CLI_VERSION + " -- HEAD")
		},
	})
}
