package types_parser

import (
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

type typeParser struct {
	schema *schemas.Schema

	enumsPkg      string
	typesPkg      string
	eventsPkg     string
	entitiesPkg   string
	repositoryPkg string
	usecasePkg    string

	typesToAvoidDuplication map[string]*Type
	enumsToAvoidDuplication map[string]*Enum

	enums      []*Enum
	types      []*Type
	events     []*Type
	entities   []*Type
	repository []*Type
	usecase    []*Type

	imports map[string]bool
}

type NewTypeParserInput struct {
	Schema        *schemas.Schema
	EnumsPkg      string
	TypesPkg      string
	EventsPkg     string
	EntitiesPkg   string
	RepositoryPkg string
	UsecasePkg    string
}

func NewTypeParser(i *NewTypeParserInput) (TypeParser, error) {
	return &typeParser{
		schema: i.Schema,

		enumsPkg:      i.EnumsPkg,
		typesPkg:      i.TypesPkg,
		eventsPkg:     i.EventsPkg,
		entitiesPkg:   i.EntitiesPkg,
		repositoryPkg: i.RepositoryPkg,
		usecasePkg:    i.UsecasePkg,

		typesToAvoidDuplication: map[string]*Type{},
		enumsToAvoidDuplication: map[string]*Enum{},
		types:                   []*Type{},
		enums:                   []*Enum{},
		imports:                 map[string]bool{},
	}, nil
}
