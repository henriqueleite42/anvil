package parser_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/anvil/anvil/internal/parser"
)

func EmailMailer() {
	schema, err := parser.ParseAnvToAnvp("../examples/advanced/email-mailer.anv")
	if err != nil {
		log.Fatal(err)
	}

	json, _ := json.Marshal(schema)

	fmt.Println(string(json))
}
