package cmd

import (
	"log"

	"github.com/henriqueleite42/anvil/cli/internal/files"
	"github.com/henriqueleite42/anvil/cli/internal/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func addParseCommand(rootCmd *cobra.Command) {
	parseCmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse the file to create the formatted version",
		Run: func(cmd *cobra.Command, args []string) {
			schema, err := parser.ParseAnvToAnvp(schemaFile)
			if err != nil {
				log.Fatal(err)
			}

			if silent {
				return
			}

			_, err = files.WriteAnvpFile(schema, schemaFile)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	parseCmd.PersistentFlags().StringVar(&schemaFile, "schema", "", "config file")
	parseCmd.MarkPersistentFlagRequired("schema")
	viper.BindPFlag("schema", parseCmd.PersistentFlags().Lookup("schema"))

	parseCmd.PersistentFlags().BoolVar(&silent, "silent", false, "if it should have an effect or only run it silently")
	viper.BindPFlag("silent", rootCmd.PersistentFlags().Lookup("silent"))

	rootCmd.AddCommand(parseCmd)
}
