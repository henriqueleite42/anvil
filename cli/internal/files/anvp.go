package files

import (
	"fmt"
	"os"
	"time"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	"gopkg.in/yaml.v3"
)

func GetAnvpFilePath(config *schemas.Config, createFolders bool) (string, error) {
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

	return fmt.Sprintf("%s/%s.anvp", path, config.ProjectName), nil
}

func GetTimestampedAnvpFilePath(config *schemas.Config, createFolders bool) (string, error) {
	myDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	path := myDir + fmt.Sprintf("/anvil/processed/%s", config.ProjectName)

	if createFolders {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	now := time.Now().UTC()
	year := fmt.Sprintf("%02d", now.Year())
	month := fmt.Sprintf("%02d", int(now.Month()))
	day := fmt.Sprintf("%02d", now.Day())
	hour := fmt.Sprintf("%02d", now.Hour())
	minute := fmt.Sprintf("%02d", now.Minute())
	seconds := fmt.Sprintf("%02d", now.Second())
	timestamp := year + month + day + hour + minute + seconds

	return fmt.Sprintf("%s/%s_%s.anvp", path, timestamp, config.ProjectName), nil
}

func WriteAnvpFile(config *schemas.Config, schema *schemas.AnvpSchema) (string, error) {
	yamlData, err := yaml.Marshal(schema)
	if err != nil {
		return "", err
	}

	filePath, err := GetAnvpFilePath(config, true)
	if err != nil {
		return "", err
	}
	timestampedFilePath, err := GetTimestampedAnvpFilePath(config, true)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(filePath, yamlData, 0644)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(timestampedFilePath, yamlData, 0644)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func ReadAnvpFile(config *schemas.Config) (*schemas.AnvpSchema, error) {
	path, err := GetAnvpFilePath(config, false)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	schema := schemas.AnvpSchema{}
	err = yaml.NewDecoder(file).Decode(&schema)
	if err != nil {
		return nil, err
	}

	return &schema, nil
}
