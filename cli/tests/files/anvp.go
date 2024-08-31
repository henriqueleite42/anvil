package files_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/henriqueleite42/anvil/cli/internal/files"
)

func ReadAnvpFile(logJson bool) {
	schema, err := files.ReadAnvpFile("../examples/advanced/authentication.anv")
	if err != nil {
		log.Fatal(err)
	}

	if logJson {
		json, _ := json.Marshal(schema)

		fmt.Println(string(json))
	}
}
