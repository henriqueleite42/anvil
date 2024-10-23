package internal

import (
	"os"

	"github.com/henriqueleite42/anvil/language-helpers/golang/formatter"
)

func WriteFile(domain string, outputFolderPath *string, fileName string, content string) error {
	myDir, err := os.Getwd()
	if err != nil {
		return err
	}

	domainKebab := formatter.PascalToKebab(domain)

	path := ""
	if outputFolderPath == nil {
		path = myDir + "/" + domainKebab
	} else {
		path = myDir + "/" + *outputFolderPath + "/" + domainKebab
	}

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path+"/"+fileName+".go", []byte(content), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
