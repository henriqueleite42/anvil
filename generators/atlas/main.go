package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/henriqueleite42/anvil/generators/atlas/internal"
	"github.com/henriqueleite42/anvil/generators/atlas/internal/postgres"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func main() {
	if len(os.Args) == 0 {
		log.Fatal("no args provided")
	}

	command := os.Args[1]
	if command != "gen" {
		log.Fatal(fmt.Sprintf("invalid command \"%s\"", command))
	}

	var schemaString string
	var outputFolderPath string
	var silent bool
	for idx, arg := range os.Args {
		if !strings.HasPrefix(arg, "--") {
			continue
		}

		if arg == "--schema" {
			schemaString = os.Args[idx+1]
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

	if schemaString == "" {
		log.Fatal("schema is required")
	}

	schema := &schemas.Schema{}
	err := json.Unmarshal([]byte(schemaString), schema)
	if err != nil {
		log.Fatal(err)
	}

	result, err := postgres.Parse(schema)
	if err != nil {
		log.Fatal(err)
	}

	if !silent {
		err := internal.WriteHclFile(outputFolderPath, schema, result)
		if err != nil {
			log.Fatal(err)
		}
	}
}
