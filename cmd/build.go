package cmd

import (
	"github.com/anvil/anvil/cmd/build"
	"github.com/anvil/anvil/cmd/parse"
	"github.com/spf13/cobra"
)

func addBuildCommand(rootCmd *cobra.Command) {
	parseCmd := &cobra.Command{
		Use:   "build",
		Short: "Build the file to check for errors",
		Run: func(cmd *cobra.Command, args []string) {
			schemaFile := cmd.Flag("schema").Value.String()

			schema := parse.Parse(schemaFile)

			build.Build(schemaFile, schema)
		},
	}

	rootCmd.AddCommand(parseCmd)
}
