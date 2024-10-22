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
		Use:     "generate",
		Aliases: []string{"gen"},
		Args:    cobra.ExactArgs(1),
		Short:   "Generate the file to check for errors",
		Run: func(cmd *cobra.Command, args []string) {
			generator := args[0]

			argsToGenerator := []string{
				"gen",
			}

			schema, err := parser.ParseAnvToAnvp(schemaFiles)
			if err != nil {
				log.Fatal(err)
			}

			anvpPath, err := files.WriteAnvpFile(schema, schemaFiles)
			if err != nil {
				log.Fatal(err)
			}

			argsToGenerator = append(argsToGenerator,
				"--schema",
				anvpPath,
			)

			if silent {
				argsToGenerator = append(argsToGenerator, "--silent")
			}

			if outputFolderPath != "" {
				argsToGenerator = append(argsToGenerator, "--outDir", outputFolderPath)
			}

			generatorPath := config.GetConfigPath() + "/generators/" + generator

			if _, err := os.Stat(generatorPath); errors.Is(err, os.ErrNotExist) {
				log.Fatalf("generator \"%s\" isn't installed. Run `anvil install %s <generator download uri>` to install it and run the command again", generator, generator)
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

	buildCmd.PersistentFlags().BoolVar(&silent, "silent", false, "if it should have an effect or only run it silently")
	viper.BindPFlag("silent", rootCmd.PersistentFlags().Lookup("silent"))

	rootCmd.AddCommand(buildCmd)
}
