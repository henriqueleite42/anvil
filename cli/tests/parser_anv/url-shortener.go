package parser_anv_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/anvil/anvil/internal/parser_anv"
)

func UrlShortener(logJson bool) {
	schema, err := parser_anv.ParseAnvToAnvp("../examples/intermediary/url-shortener.anv")
	if err != nil {
		log.Fatal(err)
	}

	if logJson {
		json, _ := json.Marshal(schema)

		fmt.Println(string(json))
	}
}
