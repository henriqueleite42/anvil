package files

import (
	"fmt"
	"os"

	"github.com/anvil/anvil/internal/formatter"
	"github.com/anvil/anvil/schemas"
	"gopkg.in/yaml.v3"
)

func WriteFile(schema *schemas.Schema) error {
	yamlData, err := yaml.Marshal(schema)
	if err != nil {
		return err
	}

	myDir, err := os.Getwd()
	if err != nil {
		return err
	}

	path := myDir + "/anvil"

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	domainKebab := formatter.PascalToKebab(schema.Domain)

	err = os.WriteFile(fmt.Sprintf("%s/%s.anvp", path, domainKebab), yamlData, 0644)
	if err != nil {
		return err
	}

	return nil
}
