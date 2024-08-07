package build

import (
	"os"
	"path/filepath"

	"github.com/anuntech/hephaestus/cmd/schema"
	"gopkg.in/yaml.v3"
)

func Build(schemaFile string, schema *schema.Schema) {
	d, err := yaml.Marshal(schema)
	if err != nil {
		panic("fail")
	}

	path := "./hephaestus/" + schemaFile

	if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
		panic(err.Error())
	}

	f, err := os.Create(path)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	_, err = f.Write(d)
	if err != nil {
		panic(err.Error())
	}
}
