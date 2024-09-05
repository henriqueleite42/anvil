package files

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func readRemoteFile(uri string) (map[string]any, error) {
	// TODO do something to accept remote repository requests and use ssh keys to authenticate, like go does
	res, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}

	fileData, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	data := map[string]any{}
	err = yaml.Unmarshal(fileData, &data)
	if err != nil {
		log.Fatal(err)
	}

	return data, nil
}

func readLocalFile(uri string) (map[string]any, error) {
	fileData, err := os.ReadFile(uri)
	if err != nil {
		log.Fatal(err)
	}

	data := map[string]any{}
	err = yaml.Unmarshal(fileData, &data)
	if err != nil {
		log.Fatal(err)
	}

	return data, nil
}

func ReadAnvFile(uri string) (map[string]any, error) {
	if strings.HasPrefix(uri, "http") {
		return readRemoteFile(uri)
	}

	return readLocalFile(uri)
}
