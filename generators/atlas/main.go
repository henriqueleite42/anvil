package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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
		log.Fatal(fmt.Sprintf("invalid command \"%s\"", command))
	}

	var schemaPath string
	var outputFolderPath string
	var silent bool
	for idx, arg := range os.Args {
		if !strings.HasPrefix(arg, "--") {
			continue
		}

		if arg == "--schema" {
			schemaPath = os.Args[idx+1]
			continue
		}

		if arg == "--outDir" {
			outputFolderPath = os.Args[idx+1]
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

	schema := &schemas.Schema{}
	err = yaml.Unmarshal(fileData, &schema)
	if err != nil {
		log.Fatal(err)
	}

	err = postgres.Parse(schema, silent, outputFolderPath)
	if err != nil {
		log.Fatal(err)
	}
}
