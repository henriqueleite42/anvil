package internal

import (
	"fmt"
	"os"

	"github.com/henriqueleite42/anvil/cli/formatter"
	"github.com/henriqueleite42/anvil/cli/schemas"
)

func WriteFile(path string, kind string, schema *schemas.Schema, content string) error {
	myDir, err := os.Getwd()
	if err != nil {
		return err
	}

	domainKebab := formatter.PascalToKebab(schema.Domain)

	if path == "" {
		path = fmt.Sprintf("%s/%s", myDir, kind)
	} else {
		path = fmt.Sprintf("%s/%s/%s", myDir, path, kind)
	}

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path+"/"+domainKebab+".go", []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}
