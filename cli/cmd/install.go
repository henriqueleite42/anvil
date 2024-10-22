package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/henriqueleite42/anvil/cli/cmd/config"
	"github.com/spf13/cobra"
)

func readRemoteFile(uri string) ([]byte, error) {
	// TODO do something to accept remote repository requests and use ssh keys to authenticate, like go does
	res, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}

	return io.ReadAll(res.Body)
}

func readLocalFile(uri string) ([]byte, error) {
	return os.ReadFile(uri)
}

func readFile(uri string) ([]byte, error) {
	if strings.HasPrefix(uri, "http") {
		return readRemoteFile(uri)
	}

	return readLocalFile(uri)
}

func addInstallCommand(rootCmd *cobra.Command) {
	buildCmd := &cobra.Command{
		Use: "install [generator name] [generator uri] [generator version]",
		Aliases: []string{
			"install",
			"add",
			"i",
		},
		Short: "Install generators for Anvil",
		Args:  cobra.MatchAll(cobra.ExactArgs(3), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			uri := args[1]
			version := args[2]

			fileData, err := readFile(uri)
			if err != nil {
				log.Fatal(err)
			}

			dir := config.GetConfigPath()
			path := fmt.Sprintf("%s/generators/%s/%s", dir, name, version)

			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}

			err = os.WriteFile(path+"/bin", fileData, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("generator \"%s\" installed successfully%s", name, "\n")
		},
	}

	rootCmd.AddCommand(buildCmd)
}
