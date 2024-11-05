package files

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	"gopkg.in/yaml.v3"
)

func ReadConfigFile(uri string) (*schemas.Config, error) {
	fileData, err := os.ReadFile(uri)
	if err != nil {
		log.Fatal(err)
	}

	data := &schemas.Config{}
	err = yaml.Unmarshal(fileData, data)
	if err != nil {
		log.Fatal(err)
	}

	vald := validator.New(validator.WithRequiredStructEnabled())

	err = vald.Struct(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
