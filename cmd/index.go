package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	schemaFile string // Used for flags.

	rootCmd = &cobra.Command{
		Use:   "hephaestus",
		Short: "God of micro-services generation",
		Long:  `hephaestus allows you to generate micro-services from schema definitions, helping you to standardize everything, avoiding human error, decreasing the learning curve and saving lots of time and money.`,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&schemaFile, "schema", "", "config file")
	viper.BindPFlag("schema", rootCmd.PersistentFlags().Lookup("schema"))

	addVersionCommand(rootCmd)
	addParseCommand(rootCmd)
	addBuildCommand(rootCmd)
}

func initConfig() {
	viper.SetConfigFile(schemaFile)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using schema file:", viper.ConfigFileUsed())
	}
}
