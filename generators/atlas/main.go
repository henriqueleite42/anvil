package main

import (
	"log"
	"os"
	"strings"

	generator_config "github.com/henriqueleite42/anvil/generators/atlas/config"
	"github.com/henriqueleite42/anvil/generators/atlas/internal/postgres"
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

	err = postgres.Parse(schema, silent, config.OutDir)
	if err != nil {
		log.Fatal(err)
	}
}
