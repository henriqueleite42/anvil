package internal

import (
	"fmt"
	"os"
	"strings"
)

func WriteFile(path string, fileNameWithPath string, content string) error {
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

	err = os.WriteFile(path+"/"+fileName, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}
