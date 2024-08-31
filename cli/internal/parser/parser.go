package parser

import "github.com/anvil/anvil/schemas"

type anvToAnvpParser struct {
	schema *schemas.Schema

	baseRef string

	filePath string
}

type resolveInput struct {
	path string
	ref  string
	k    string
	v    any
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

	err = self.resolveEntitiesMetadata(file)
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

	err = self.entities(file)
	if err != nil {
		return err
	}

	err = self.repository(file)
	if err != nil {
		return err
	}

	return nil
}

func ParseAnvToAnvp(uri string) (*schemas.Schema, error) {
	parser := &anvToAnvpParser{
		schema:   &schemas.Schema{},
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
