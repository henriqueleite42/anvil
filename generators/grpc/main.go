package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/henriqueleite42/anvil/cli/schemas"
	"github.com/henriqueleite42/anvil/generators/grpc/internal"
	"golang.org/x/exp/slices"
)

func main() {
	schemaString := os.Args[1]
	if schemaString == "" {
		log.Fatal("schema is required")
	}

	schema := &schemas.Schema{}
	err := json.Unmarshal([]byte(schemaString), schema)
	if err != nil {
		log.Fatal(err)
	}

	result, err := internal.Parse(schema)
	if err != nil {
		log.Fatal(err)
	}

	remainingArgs := os.Args[2:]

	if !slices.Contains(remainingArgs, "--silent") {
		err := internal.WriteProtoFile(schema, result)
		if err != nil {
			log.Fatal(err)
		}
	}
}