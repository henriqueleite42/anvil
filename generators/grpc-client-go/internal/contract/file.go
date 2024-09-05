package contract

import (
	"os"

	"github.com/henriqueleite42/anvil/cli/formatter"
	"github.com/henriqueleite42/anvil/cli/schemas"
)

func WriteContractFile(path string, schema *schemas.Schema, content string) error {
	myDir, err := os.Getwd()
	if err != nil {
		return err
	}

	domainKebab := formatter.PascalToKebab(schema.Domain)

	if path == "" {
		path = myDir + "/" + domainKebab
	} else {
		path = myDir + "/" + path + "/" + domainKebab
	}

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path+"/contract.go", []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}
