package types_parser

import (
	"github.com/henriqueleite42/anvil/language-helpers/golang/imports"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

type GetEnumImport func(e *schemas.Enum) *imports.Import
type GetTypeImport func(t *schemas.Type) *imports.Import

type typeParser struct {
	schema *schemas.AnvpSchema

	getEnumsImport      GetEnumImport
	getTypesImport      GetTypeImport
	getEventsImport     GetTypeImport
	getEntitiesImport   GetTypeImport
	getRepositoryImport GetTypeImport
	getUsecaseImport    GetTypeImport

	typesToAvoidDuplication map[string]*Type
	enumsToAvoidDuplication map[string]*Enum

	enums      []*Enum
	types      []*Type
	events     []*Type
	entities   []*Type
	repository []*Type
	usecase    []*Type
}

type NewTypeParserInput struct {
	Schema *schemas.AnvpSchema

	GetEnumsImport      GetEnumImport
	GetTypesImport      GetTypeImport
	GetEventsImport     GetTypeImport
	GetEntitiesImport   GetTypeImport
	GetRepositoryImport GetTypeImport
	GetUsecaseImport    GetTypeImport
}

func NewTypeParser(i *NewTypeParserInput) (TypesParser, error) {
	return &typeParser{
		schema: i.Schema,

		getEnumsImport:      i.GetEnumsImport,
		getTypesImport:      i.GetTypesImport,
		getEventsImport:     i.GetEventsImport,
		getEntitiesImport:   i.GetEntitiesImport,
		getRepositoryImport: i.GetRepositoryImport,
		getUsecaseImport:    i.GetUsecaseImport,

		typesToAvoidDuplication: map[string]*Type{},
		enumsToAvoidDuplication: map[string]*Enum{},

		types:      []*Type{},
		enums:      []*Enum{},
		events:     []*Type{},
		entities:   []*Type{},
		repository: []*Type{},
		usecase:    []*Type{},
	}, nil
}
