package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/henriqueleite42/anvil/cli/schemas"
	"github.com/henriqueleite42/anvil/generators/atlas/internal"
	"github.com/henriqueleite42/anvil/generators/atlas/internal/postgres"
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
		fmt.Println(err.Error())
		log.Fatal(err)
	}

	result, err := postgres.Parse(schema)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

	remainingArgs := os.Args[2:]

	if !slices.Contains(remainingArgs, "--silent") {
		err := internal.WriteHclFile(schema, result)
		if err != nil {
			fmt.Println(err.Error())
			log.Fatal(err)
		}
	}
}
