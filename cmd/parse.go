package cmd

import (
	"github.com/anvil/anvil/cmd/parse"
	"github.com/spf13/cobra"
)

func addParseCommand(rootCmd *cobra.Command) {
	parseCmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse the file to check for errors",
		Run: func(cmd *cobra.Command, args []string) {
			schemaFile := cmd.Flag("schema").Value.String()

			parse.Parse(schemaFile)
		},
	}

	rootCmd.AddCommand(parseCmd)
}
