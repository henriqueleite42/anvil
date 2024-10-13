package postgres

import (
	"os"

	"github.com/henriqueleite42/anvil/language-helpers/golang/formatter"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func WriteFile(schema *schemas.Schema, outputFolderPath string, content string) error {
	myDir, err := os.Getwd()
	if err != nil {
		return err
	}

	domainKebab := formatter.PascalToKebab(schema.Domain)

	path := myDir
	if outputFolderPath != "" {
		path = myDir + "/" + outputFolderPath
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
