package types_parser

import (
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

type typeParser struct {
	schema *schemas.AnvpSchema

	moduleName string

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

	importsTypes      map[string]bool
	importsEnums      map[string]bool
	importsEvents     map[string]bool
	importsEntities   map[string]bool
	importsRepository map[string]bool
	importsUsecase    map[string]bool
}

type NewTypeParserInput struct {
	Schema        *schemas.AnvpSchema
	ModuleName    string
	EnumsPkg      string
	TypesPkg      string
	EventsPkg     string
	EntitiesPkg   string
	RepositoryPkg string
	UsecasePkg    string
}

func NewTypeParser(i *NewTypeParserInput) (TypesParser, error) {
	return &typeParser{
		schema: i.Schema,

		moduleName: i.ModuleName,

		enumsPkg:      i.EnumsPkg,
		typesPkg:      i.TypesPkg,
		eventsPkg:     i.EventsPkg,
		entitiesPkg:   i.EntitiesPkg,
		repositoryPkg: i.RepositoryPkg,
		usecasePkg:    i.UsecasePkg,

		typesToAvoidDuplication: map[string]*Type{},
		enumsToAvoidDuplication: map[string]*Enum{},

		types:      []*Type{},
		enums:      []*Enum{},
		events:     []*Type{},
		entities:   []*Type{},
		repository: []*Type{},
		usecase:    []*Type{},

		importsTypes:      map[string]bool{},
		importsEnums:      map[string]bool{},
		importsEvents:     map[string]bool{},
		importsEntities:   map[string]bool{},
		importsRepository: map[string]bool{},
		importsUsecase:    map[string]bool{},
	}, nil
}
