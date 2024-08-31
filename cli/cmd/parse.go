package cmd

import (
	"log"

	"github.com/anvil/anvil/internal/files"
	"github.com/anvil/anvil/internal/parser_anv"
	"github.com/spf13/cobra"
)

func addParseCommand(rootCmd *cobra.Command) {
	parseCmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse the file to create the formatted version",
		Run: func(cmd *cobra.Command, args []string) {
			schemaFile := cmd.Flag("schema").Value.String()

			schema, err := parser_anv.ParseAnvToAnvp(schemaFile)
			if err != nil {
				log.Fatal(err)
			}

			err = files.WriteFile(schema)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	rootCmd.AddCommand(parseCmd)
}
