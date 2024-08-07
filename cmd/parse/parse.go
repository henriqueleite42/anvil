package parse

import (
	"log"
	"os"

	"github.com/anuntech/hephaestus/cmd/schema"
	"gopkg.in/yaml.v3"
)

func Parse(schemaFile string) *schema.Schema {
	fileData, err := os.ReadFile(schemaFile)
	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string]any)
	err = yaml.Unmarshal(fileData, &data)
	if err != nil {
		log.Fatal(err)
	}

	schema, err := file(data)
	if err != nil {
		log.Fatal(err)
	}

	return schema
}
