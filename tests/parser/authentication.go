package parser_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/anvil/anvil/internal/parser"
	"gopkg.in/yaml.v3"
)

func Authentication() {
	fileData, err := os.ReadFile("./examples/advanced/authentication.anv")
	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string]any)
	err = yaml.Unmarshal(fileData, &data)
	if err != nil {
		log.Fatal(err)
	}

	p := parser.NewParser()

	err = p.Parse(data)
	if err != nil {
		log.Fatal(err)
	}

	s := p.GetSchema()

	json, _ := json.Marshal(s)

	fmt.Println(string(json))
}
