package internal

import (
	"os"

	"github.com/henriqueleite42/anvil/cli/formatter"
)

func WriteFile(domain string, path string, fileName string, content string) error {
	myDir, err := os.Getwd()
	if err != nil {
		return err
	}

	domainKebab := formatter.PascalToKebab(domain)

	if path == "" {
		path = myDir + "/" + domainKebab
	} else {
		path = myDir + "/" + path + "/" + domainKebab
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
