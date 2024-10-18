package files

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/henriqueleite42/anvil/language-helpers/golang/hashing"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	"gopkg.in/yaml.v3"
)

func GetAnvpFilePath(anvpFileName string, createFolders bool) (string, error) {
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

	return fmt.Sprintf("%s/%s.anvp", path, anvpFileName), nil
}

func WriteAnvpFile(schema *schemas.AnvpSchema, schemaFiles []string) (string, error) {
	yamlData, err := yaml.Marshal(schema)
	if err != nil {
		return "", err
	}

	sort.Slice(schemaFiles, func(i, j int) bool {
		return schemaFiles[i] < schemaFiles[j]
	})

	anvpFileName := hashing.String(
		strings.Join(schemaFiles, ""),
	)

	filePath, err := GetAnvpFilePath(anvpFileName, true)
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
