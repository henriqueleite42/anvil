package parser

import (
	"github.com/henriqueleite42/anvil/cli/internal/files"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

type anvToAnvpParser struct {
	schema *schemas.Schema

	filePath string
}

type resolveInput struct {
	namePrefix string // Internal use. Correctly parse child map types with the prefix of their parent.

	path string // Original path
	ref  string // Ref until now
	k    string // Key being resolved, usually the type name, but if it's an child type, it's only part of the type's name
	v    any    // Value, type specification
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

	err = self.usecase(file)
	if err != nil {
		return err
	}

	err = self.delivery(file)
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

	file, err := files.ReadAnvFile(uri)
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
