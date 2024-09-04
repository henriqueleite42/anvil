package files

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func ReadAnvFile(uri string) (map[string]any, error) {
	fileData, err := os.ReadFile(uri)
	if err != nil {
		log.Fatal(err)
	}

	data := map[string]any{}
	err = yaml.Unmarshal(fileData, &data)
	if err != nil {
		log.Fatal(err)
	}

	return data, nil
}
