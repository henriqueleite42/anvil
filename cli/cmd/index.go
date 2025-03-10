package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	global_Silent bool

	rootCmd = &cobra.Command{
		Use:   "anvil",
		Short: "Generate code from a common schema",
		Long:  `anvil allows you to generate micro-services from schema definitions, helping you to standardize everything, avoiding human error, decreasing the learning curve and saving lots of time and money.`,
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

	addVersionCommand(rootCmd)
	addParseCommand(rootCmd)
	addGenerateCommand(rootCmd)
	addInstallCommand(rootCmd)
}

func initConfig() {
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
