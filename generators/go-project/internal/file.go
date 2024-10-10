package internal

import (
	"fmt"
	"os"
	"strings"
)

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func WriteFile(path string, fileNameWithPath string, content string, overwrite bool) error {
	myDir, err := os.Getwd()
	if err != nil {
		return err
	}

	partsToFile := strings.Split(fileNameWithPath, "/")
	fileName := partsToFile[len(partsToFile)-1]
	partsToFile = partsToFile[:len(partsToFile)-1]
	pathToFile := strings.Join(partsToFile, "/")

	if path == "" {
		path = fmt.Sprintf("%s/%s", myDir, pathToFile)
	} else {
		path = fmt.Sprintf("%s/%s/%s", myDir, path, pathToFile)
	}

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	finalPath := path + "/" + fileName

	if !overwrite && fileExists(finalPath) {
		return fmt.Errorf("file \"%s\" already exists, generator \"go-project\" is unable to overwrite files without losing data", finalPath)
	}

	err = os.WriteFile(finalPath, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}
