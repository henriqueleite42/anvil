package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/anuntech/hephaestus/cmd/parse"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func addParseCommand(rootCmd *cobra.Command) {
	parseCmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse the file to check for errors",
		Run: func(cmd *cobra.Command, args []string) {
			schemaFile := cmd.Flag("schema").Value.String()

			fileData, err := os.ReadFile(schemaFile)
			if err != nil {
				log.Fatal(err)
			}

			data := make(map[string]any)
			err = yaml.Unmarshal(fileData, &data)
			if err != nil {
				log.Fatal(err)
			}

			schema, err := parse.File(data)
			if err != nil {
				log.Fatal(err)
			}

			d, err := json.Marshal(schema)
			if err != nil {
				panic("fail")
			}
			fmt.Println(string(d))
		},
	}

	rootCmd.AddCommand(parseCmd)
}
