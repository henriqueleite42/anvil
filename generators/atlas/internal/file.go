package internal

import (
	"os"

	"github.com/henriqueleite42/anvil/cli/schemas"
)

func WriteHclFile(schema *schemas.Schema, content string) error {
	myDir, err := os.Getwd()
	if err != nil {
		return err
	}

	path := myDir + "/" + schema.Domain + ".hcl"

	err = os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}
