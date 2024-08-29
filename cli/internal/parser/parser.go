package parser

import (
	"github.com/anvil/anvil/internal/schema"
)

type anvToAnvpParser struct {
	schema *schema.Schema

	basePath string
	filePath string
}

func (self *anvToAnvpParser) parse(file map[string]any) error {
	err := self.domain(file)
	if err != nil {
		return err
	}

	err = self.metadata(file)
	if err != nil {
		return err
	}

	err = self.relationships(file)
	if err != nil {
		return err
	}

	err = self.imports(file)
	if err != nil {
		return err
	}

	err = self.auth(file)
	if err != nil {
		return err
	}

	err = self.enums(file)
	if err != nil {
		return err
	}

	err = self.types(file)
	if err != nil {
		return err
	}

	err = self.events(file)
	if err != nil {
		return err
	}

	return nil
}

func ParseAnvToAnvp(uri string) (*schema.Schema, error) {
	parser := &anvToAnvpParser{
		schema:   &schema.Schema{},
		filePath: uri,
	}

	file, err := parser.readAnvFile(uri)
	if err != nil {
		return nil, err
	}

	err = parser.parse(file)
	if err != nil {
		return nil, err
	}

	err = parser.stateHashes()
	if err != nil {
		return nil, err
	}

	return parser.schema, nil
}
