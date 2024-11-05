package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/henriqueleite42/anvil/cli/cmd/config"
	"github.com/henriqueleite42/anvil/cli/internal/files"
	"github.com/henriqueleite42/anvil/cli/internal/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func addGenerateCommand(rootCmd *cobra.Command) {
	generateCmd := &cobra.Command{
		Use: "generate [config file path]",
		Aliases: []string{
			"generate",
			"gen",
		},
		Args:  cobra.MaximumNArgs(1),
		Short: "Run generators based on a config file",
		Run: func(cmd *cobra.Command, args []string) {
			var configFilePath string
			if len(args) > 0 {
				configFilePath = args[0]
			}
			if configFilePath == "" {
				configFilePath = "./anvil.yaml"
			}

			configFile, err := files.ReadConfigFile(configFilePath)
			if err != nil {
				log.Fatal(err)
			}

			argsToGenerator := []string{
				"gen",
			}

			schema, err := parser.ParseAnvToAnvp(configFile.Schemas)
			if err != nil {
				log.Fatal(err)
			}

			anvpPath, err := files.WriteAnvpFile(configFile, schema)
			if err != nil {
				log.Fatal(err)
			}

			argsToGenerator = append(argsToGenerator,
				"--config",
				configFilePath,
			)
			argsToGenerator = append(argsToGenerator,
				"--schema",
				anvpPath,
			)

			if global_Silent {
				argsToGenerator = append(argsToGenerator, "--silent")
			}

			for _, v := range configFile.Generators {
				generatorPath := fmt.Sprintf(
					"%s/generators/%s/%s/bin",
					config.GetConfigPath(),
					v.Name,
					v.Version,
				)

				if _, err := os.Stat(generatorPath); errors.Is(err, os.ErrNotExist) {
					log.Fatalf("generator \"%s\" isn't installed. Run `anvil i %s <generator download uri> %s` to install it and run the command again", v.Name, v.Name, v.Version)
				}

				generatorCmd := exec.Command(generatorPath, argsToGenerator...)
				generatorCmd.Stdout = os.Stdout
				generatorCmd.Stderr = os.Stderr
				err = generatorCmd.Run()
				if err != nil {
					log.Fatal(err)
				}
			}
		},
	}

	generateCmd.PersistentFlags().BoolVar(&global_Silent, "silent", false, "if it should have an effect or only run it silently")
	viper.BindPFlag("silent", rootCmd.PersistentFlags().Lookup("silent"))

	rootCmd.AddCommand(generateCmd)
}
