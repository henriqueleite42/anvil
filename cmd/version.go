package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func addVersionCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number of hephaestus",
		Long:  `All software has versions. This is hephaestus's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hephaestus Micro-Services Generator v0.1 -- HEAD")
		},
	})
}
