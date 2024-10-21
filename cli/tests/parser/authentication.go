package parser_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/henriqueleite42/anvil/cli/internal/parser"
)

func Authentication(logJson bool) {
	schema, err := parser.ParseAnvToAnvp(
		[]string{
			"../examples/advanced/authentication.anv",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	if logJson {
		json, _ := json.Marshal(schema)

		fmt.Println(string(json))
	}
}
