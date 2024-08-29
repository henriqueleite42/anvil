package cmd

import (
	"log"

	"github.com/anvil/anvil/internal/parser"
	"github.com/spf13/cobra"
)

func addBuildCommand(rootCmd *cobra.Command) {
	parseCmd := &cobra.Command{
		Use:   "build",
		Short: "Build the file to check for errors",
		Run: func(cmd *cobra.Command, args []string) {
			schemaFile := cmd.Flag("schema").Value.String()

			_, err := parser.ParseAnvToAnvp(schemaFile)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	rootCmd.AddCommand(parseCmd)
}
