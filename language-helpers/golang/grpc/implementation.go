package grpc

import (
	"github.com/henriqueleite42/anvil/language-helpers/golang/grpc/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	"github.com/henriqueleite42/anvil/language-helpers/golang/template"
	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

type goGrpcParser struct {
	schema *schemas.AnvpSchema

	imports map[string]bool

	enumsPkg      string
	typesPkg      string
	eventsPkg     string
	entitiesPkg   string
	repositoryPkg string
	usecasePkg    string

	goTypeParser    types_parser.TypesParser
	templateManager template.TemplateManager
}

type NewGrpcParserInput struct {
	Schema *schemas.AnvpSchema

	ModuleName string

	EnumsPkg      string
	TypesPkg      string
	EventsPkg     string
	EntitiesPkg   string
	RepositoryPkg string
	UsecasePkg    string

	GoTypeParser types_parser.TypesParser
}

func NewGrpcParser(i *NewGrpcParserInput) GrpcParser {
	templateManager := template.NewTemplateManager()

	templateManager.AddTemplate("input-prop-list", templates.InputPropListTempl)
	templateManager.AddTemplate("input-prop-map", templates.InputPropMapTempl)
	templateManager.AddTemplate("input-prop-optional", templates.InputPropOptionalTempl)

	return &goGrpcParser{
		schema: i.Schema,

		imports: map[string]bool{},

		goTypeParser:    i.GoTypeParser,
		templateManager: templateManager,
	}
}
