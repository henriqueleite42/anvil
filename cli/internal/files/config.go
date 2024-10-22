package files

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	"gopkg.in/yaml.v3"
)

func readConfigRemoteFile(uri string) (*schemas.Config, error) {
	// TODO do something to accept remote repository requests and use ssh keys to authenticate, like go does
	res, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}

	fileData, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	data := &schemas.Config{}
	err = yaml.Unmarshal(fileData, data)
	if err != nil {
		log.Fatal(err)
	}

	return data, nil
}

func readConfigLocalFile(uri string) (*schemas.Config, error) {
	fileData, err := os.ReadFile(uri)
	if err != nil {
		log.Fatal(err)
	}

	data := &schemas.Config{}
	err = yaml.Unmarshal(fileData, data)
	if err != nil {
		log.Fatal(err)
	}

	return data, nil
}

func ReadConfigFile(uri string) (*schemas.Config, error) {
	if strings.HasPrefix(uri, "http") {
		return readConfigRemoteFile(uri)
	}

	return readConfigLocalFile(uri)
}
