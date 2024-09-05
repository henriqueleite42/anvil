package internal

import (
	"os"

	"github.com/henriqueleite42/anvil/cli/formatter"
	"github.com/henriqueleite42/anvil/cli/schemas"
)

func WriteHclFile(path string, schema *schemas.Schema, content string) error {
	myDir, err := os.Getwd()
	if err != nil {
		return err
	}

	domainKebab := formatter.PascalToKebab(schema.Domain)

	if path == "" {
		path = myDir
	} else {
		path = myDir + "/" + path
	}

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path+"/"+domainKebab+".hcl", []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}
