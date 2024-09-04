package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/henriqueleite42/anvil/cli/internal/files"
	"github.com/henriqueleite42/anvil/cli/internal/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func addBuildCommand(rootCmd *cobra.Command) {
	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Build the file to check for errors",
		Run: func(cmd *cobra.Command, args []string) {
			schema, err := parser.ParseAnvToAnvp(schemaFile)
			if err != nil {
				log.Fatal(err)
			}

			err = files.WriteAnvpFile(schema, schemaFile)
			if err != nil {
				log.Fatal(err)
			}

			json, err := json.Marshal(schema)
			if err != nil {
				log.Fatal(err)
			}
			jsonString := string(json)

			argsToGenerator := []string{
				"gen",
				"--schema",
				jsonString,
			}

			if silent {
				argsToGenerator = append(argsToGenerator, "--silent")
			}

			if outputFolderPath != "" {
				argsToGenerator = append(argsToGenerator, "--outDir", outputFolderPath)
			}

			for _, v := range generators {
				stdout, err := exec.Command(v, argsToGenerator...).Output()
				fmt.Println(string(stdout))
				if err != nil {
					fmt.Println(err)
				}
			}
		},
	}

	buildCmd.PersistentFlags().BoolVar(&silent, "silent", false, "if it should have an effect or only run it silently")
	viper.BindPFlag("silent", rootCmd.PersistentFlags().Lookup("silent"))

	buildCmd.PersistentFlags().StringArrayVar(&generators, "generators", []string{}, "generator to be used, can be passed more than once")
	buildCmd.MarkPersistentFlagRequired("generators")
	viper.BindPFlag("generators", rootCmd.PersistentFlags().Lookup("generators"))

	buildCmd.PersistentFlags().StringVar(&outputFolderPath, "outDir", "", "output directory path")
	viper.BindPFlag("outDir", rootCmd.PersistentFlags().Lookup("outDir"))

	rootCmd.AddCommand(buildCmd)
}
