package files_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/henriqueleite42/anvil/cli/internal/files"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func ReadAnvpFile(logJson bool) {
	schema, err := files.ReadAnvpFile(&schemas.Config{
		ProjectName: "Foo",
		Schemas: []string{
			"../examples/advanced/authentication.anv",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	if logJson {
		json, _ := json.Marshal(schema)

		fmt.Println(string(json))
	}
}
