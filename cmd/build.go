package cmd

import (
	"log"
	"os"

	"github.com/anvil/anvil/internal/parser"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func addBuildCommand(rootCmd *cobra.Command) {
	parseCmd := &cobra.Command{
		Use:   "build",
		Short: "Build the file to check for errors",
		Run: func(cmd *cobra.Command, args []string) {
			schemaFile := cmd.Flag("schema").Value.String()

			schemaParser := parser.NewParser()

			fileData, err := os.ReadFile(schemaFile)
			if err != nil {
				log.Fatal(err)
			}

			data := make(map[string]any)
			err = yaml.Unmarshal(fileData, &data)
			if err != nil {
				log.Fatal(err)
			}

			err = schemaParser.Parse(data)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	rootCmd.AddCommand(parseCmd)
}
