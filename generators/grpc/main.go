package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/henriqueleite42/anvil/cli/schemas"
	"github.com/henriqueleite42/anvil/generators/grpc/internal"
	"golang.org/x/exp/slices"
)

func main() {
	schemaString := os.Args[1]
	if schemaString == "" {
		fmt.Println("foo")
		log.Fatal("schema is required")
	}

	schema := &schemas.Schema{}
	err := json.Unmarshal([]byte(schemaString), schema)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

	result, err := internal.Parse(schema)
	if err != nil {
		fmt.Println(err.Error())
	}

	remainingArgs := os.Args[2:]

	if !slices.Contains(remainingArgs, "--silent") {
		internal.WriteProtoFile(schema, result)
	}
}
