package cmd

import (
	"log"

	"github.com/henriqueleite42/anvil/cli/internal/files"
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
		Short: "Parse the file to create the formatted version",
		Run: func(cmd *cobra.Command, args []string) {
			schema, err := parser.ParseAnvToAnvp(parse_SchemaFiles)
			if err != nil {
				log.Fatal(err)
			}

			if global_Silent {
				return
			}

			_, err = files.WriteAnvpFile(schema, parse_SchemaFiles)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	parseCmd.PersistentFlags().StringSliceVar(&parse_SchemaFiles, "schema", []string{}, "config files")
	parseCmd.MarkPersistentFlagRequired("schema")
	viper.BindPFlag("schema", parseCmd.PersistentFlags().Lookup("schema"))

	parseCmd.PersistentFlags().BoolVar(&global_Silent, "silent", false, "if it should have an effect or only run it silently")
	viper.BindPFlag("silent", rootCmd.PersistentFlags().Lookup("silent"))

	rootCmd.AddCommand(parseCmd)
}
