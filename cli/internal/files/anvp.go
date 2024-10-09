package files

import (
	"fmt"
	"os"
	"strings"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	"gopkg.in/yaml.v3"
)

func GetAnvpFilePath(anvFilePath string, createFolders bool) (string, error) {
	myDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	path := myDir + "/anvil/processed"

	if createFolders {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	parts := strings.Split(anvFilePath, "/")
	fileName := parts[len(parts)-1] + "p"

	return fmt.Sprintf("%s/%s", path, fileName), nil
}

func WriteAnvpFile(schema *schemas.Schema, anvFilePath string) (string, error) {
	yamlData, err := yaml.Marshal(schema)
	if err != nil {
		return "", err
	}

	filePath, err := GetAnvpFilePath(anvFilePath, true)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(filePath, yamlData, 0644)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func ReadAnvpFile(anvFilePath string) (*schemas.Schema, error) {
	path, err := GetAnvpFilePath(anvFilePath, false)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	schema := schemas.Schema{}
	err = yaml.NewDecoder(file).Decode(&schema)
	if err != nil {
		return nil, err
	}

	return &schema, nil
}
