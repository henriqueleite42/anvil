package parser_anv_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/anvil/anvil/internal/parser_anv"
)

func Authentication(logJson bool) {
	schema, err := parser_anv.ParseAnvToAnvp("../examples/advanced/authentication.anv")
	if err != nil {
		log.Fatal(err)
	}

	if logJson {
		json, _ := json.Marshal(schema)

		fmt.Println(string(json))
	}
}
