package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/henriqueleite42/anvil/generators/go-project/internal"
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

	files, err := internal.Parse(schema)
	if err != nil {
		log.Fatal(err)
	}

	if silent {
		return
	}

	for _, v := range files {
		err = internal.WriteFile(outputFolderPath, v.Name, v.Content)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				slog.Warn(err.Error())
			} else {
				log.Fatal(err)
			}
		}
	}
}
