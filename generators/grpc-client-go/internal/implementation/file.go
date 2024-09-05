package implementation

import (
	"os"

	"github.com/henriqueleite42/anvil/cli/formatter"
	"github.com/henriqueleite42/anvil/cli/schemas"
)

func WriteImplementationFile(path string, schema *schemas.Schema, content string) error {
	myDir, err := os.Getwd()
	if err != nil {
		return err
	}

	domainKebab := formatter.PascalToKebab(schema.Domain)

	fileWithDir := "/" + domainKebab + "/implementation.go"

	if path == "" {
		path = myDir + fileWithDir
	} else {
		path = myDir + "/" + path + fileWithDir
	}

	err = os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}
