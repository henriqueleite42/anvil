package cmd

import (
	"log"

	"github.com/henriqueleite42/anvil/cli/internal/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	parse_SchemaFiles []string
)

func addParseCommand(rootCmd *cobra.Command) {
	parseCmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse the schemas to check for errors",
		Run: func(cmd *cobra.Command, args []string) {
			_, err := parser.ParseAnvToAnvp(parse_SchemaFiles)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	parseCmd.PersistentFlags().StringSliceVar(&parse_SchemaFiles, "schema", []string{}, "config files")
	parseCmd.MarkPersistentFlagRequired("schema")
	viper.BindPFlag("schema", parseCmd.PersistentFlags().Lookup("schema"))

	rootCmd.AddCommand(parseCmd)
}
