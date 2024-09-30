package types_parser

import (
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

type typeParser struct {
	schema *schemas.Schema

	typesToAvoidDuplication map[string]*Type
	enumsToAvoidDuplication map[string]*Enum
	types                   []*Type
	enums                   []*Enum
	imports                 map[string]bool
}

func NewTypeParser(schema *schemas.Schema) (TypeParser, error) {
	return &typeParser{
		schema: schema,

		typesToAvoidDuplication: map[string]*Type{},
		enumsToAvoidDuplication: map[string]*Enum{},
		types:                   []*Type{},
		enums:                   []*Enum{},
		imports:                 map[string]bool{},
	}, nil
}
