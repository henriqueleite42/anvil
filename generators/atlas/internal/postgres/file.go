package postgres

import (
	"os"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func WriteFile(schema *schemas.AnvpSchema, outputFolderPath *string, content string) error {
	myDir, err := os.Getwd()
	if err != nil {
		return err
	}

	path := myDir
	if outputFolderPath != nil {
		path = myDir + "/" + *outputFolderPath
	}

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path+"/database.hcl", []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}
