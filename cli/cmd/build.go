package cmd

import (
	"errors"
	"log"
	"os"
	"os/exec"

	"github.com/henriqueleite42/anvil/cli/cmd/config"
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

			anvpPath, err := files.WriteAnvpFile(schema, schemaFile)
			if err != nil {
				log.Fatal(err)
			}

			argsToGenerator := []string{
				"gen",
				"--schema",
				anvpPath,
			}

			if silent {
				argsToGenerator = append(argsToGenerator, "--silent")
			}

			if outputFolderPath != "" {
				argsToGenerator = append(argsToGenerator, "--outDir", outputFolderPath)
			}

			generatorPath := config.GetConfigPath() + "/generators/" + generator

			if _, err := os.Stat(generatorPath); errors.Is(err, os.ErrNotExist) {
				log.Fatalf("generator \"%s\" isn't installed. Run `anvil install %s <generator download uri>` to install it and run the command again.", generator, generator)
			}

			generatorCmd := exec.Command(generatorPath, argsToGenerator...)
			generatorCmd.Stdout = os.Stdout
			generatorCmd.Stderr = os.Stderr
			err = generatorCmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	buildCmd.PersistentFlags().StringVar(&schemaFile, "schema", "", "config file")
	buildCmd.MarkPersistentFlagRequired("schema")
	viper.BindPFlag("schema", buildCmd.PersistentFlags().Lookup("schema"))

	buildCmd.PersistentFlags().BoolVar(&silent, "silent", false, "if it should have an effect or only run it silently")
	viper.BindPFlag("silent", rootCmd.PersistentFlags().Lookup("silent"))

	buildCmd.PersistentFlags().StringVar(&generator, "generator", "", "generator to be used")
	buildCmd.MarkPersistentFlagRequired("generator")
	viper.BindPFlag("generator", rootCmd.PersistentFlags().Lookup("generator"))

	buildCmd.PersistentFlags().StringVar(&outputFolderPath, "outDir", "", "output directory path")
	viper.BindPFlag("outDir", rootCmd.PersistentFlags().Lookup("outDir"))

	rootCmd.AddCommand(buildCmd)
}
