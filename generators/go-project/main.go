package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/henriqueleite42/anvil/cli/schemas"
	"github.com/henriqueleite42/anvil/generators/go-project/internal"
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

	models, repository, usecase, err := internal.Parse(schema)
	if err != nil {
		log.Fatal(err)
	}

	if silent {
		return
	}

	if models != "" {
		err = internal.WriteFile(outputFolderPath, "models", schema, models)
		if err != nil {
			log.Fatal(err)
		}
	}

	if repository != "" {
		err = internal.WriteFile(outputFolderPath, "repository", schema, repository)
		if err != nil {
			log.Fatal(err)
		}
	}

	if usecase != "" {
		err = internal.WriteFile(outputFolderPath, "usecase", schema, usecase)
		if err != nil {
			log.Fatal(err)
		}
	}
}