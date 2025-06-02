package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	generator_config "github.com/henriqueleite42/anvil/generators/go-project/config"
	"github.com/henriqueleite42/anvil/generators/go-project/internal"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	"gopkg.in/yaml.v3"
)

func main() {
	if len(os.Args) == 0 {
		log.Fatal("no args provided")
	}

	command := os.Args[1]
	if command != "gen" {
		log.Fatalf("invalid command \"%s\"", command)
	}

	var schemaPath string
	var configPath string
	var silent bool
	for idx, arg := range os.Args {
		if !strings.HasPrefix(arg, "--") {
			continue
		}

		if arg == "--schema" {
			schemaPath = os.Args[idx+1]
			continue
		}

		if arg == "--config" {
			configPath = os.Args[idx+1]
			continue
		}

		if arg == "--silent" {
			silent = true
			continue
		}
	}

	if schemaPath == "" {
		log.Fatal("schema is required")
	}

	fileData, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatal(err)
	}

	schema := &schemas.AnvpSchema{}
	err = yaml.Unmarshal(fileData, &schema)
	if err != nil {
		log.Fatal(err)
	}

	config := generator_config.GetConfig(configPath)

	files, err := internal.Parse(schema, config)
	if err != nil {
		log.Fatal(err)
	}

	if silent {
		return
	}

	createdFilesList := []string{}
	createdFilesCount := 0
	notModifiedFilesCount := 0
	errorFilesCount := 0
	for _, v := range files {
		err = internal.WriteFile(config.OutDir, v.Name, v.Content, v.Overwrite)
		if err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				errorFilesCount++
				slog.Error(err.Error())
			} else {
				notModifiedFilesCount++
			}
		} else {
			createdFilesCount++
			createdFilesList = append(createdFilesList, v.Name)
		}
	}

	msg := fmt.Sprintf(
		"\nAnvil generation complete!\n\nCreated: %d\n%s\n\nNot modified: %d\n\nError: %d\n",
		createdFilesCount,
		strings.Join(createdFilesList, "\n"),
		notModifiedFilesCount,
		errorFilesCount,
	)

	slog.Info(msg)
}
