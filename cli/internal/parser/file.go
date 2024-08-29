package parser

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func (self *anvToAnvpParser) readAnvFile(uri string) (map[string]any, error) {
	fileData, err := os.ReadFile(uri)
	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string]any)
	err = yaml.Unmarshal(fileData, &data)
	if err != nil {
		log.Fatal(err)
	}

	return data, nil
}
