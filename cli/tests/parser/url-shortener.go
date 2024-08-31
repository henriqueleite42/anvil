package parser_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/anvil/anvil/internal/parser"
)

func UrlShortener(logJson bool) {
	schema, err := parser.ParseAnvToAnvp("../examples/intermediary/url-shortener.anv")
	if err != nil {
		log.Fatal(err)
	}

	if logJson {
		json, _ := json.Marshal(schema)

		fmt.Println(string(json))
	}
}
