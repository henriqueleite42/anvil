package internal

import (
	"os"

	"github.com/henriqueleite42/anvil/language-helpers/golang/formatter"
)

func WriteFile(curDomain string, outputFolderPath string, content string) error {
	myDir, err := os.Getwd()
	if err != nil {
		return err
	}

	domainKebab := formatter.PascalToKebab(curDomain)

	path := myDir
	if outputFolderPath != "" {
		path = myDir + "/" + outputFolderPath
	}

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path+"/"+domainKebab+".proto", []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}
