package grpc

import (
	"github.com/henriqueleite42/anvil/language-helpers/golang/grpc/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	"github.com/henriqueleite42/anvil/language-helpers/golang/template"
	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

type goGrpcParser struct {
	schema *schemas.Schema

	goTypeParser    types_parser.TypeParser
	templateManager template.TemplateManager
}

type NewGrpcParserInput struct {
	Schema *schemas.Schema

	GoTypeParser types_parser.TypeParser
}

func NewGrpcParser(i *NewGrpcParserInput) GrpcParser {
	templateManager := template.NewTemplateManager()

	templateManager.AddTemplate("input-prop-list", templates.InputPropListTempl)
	templateManager.AddTemplate("input-prop-map", templates.InputPropMapTempl)
	templateManager.AddTemplate("input-prop-optional", templates.InputPropOptionalTempl)

	return &goGrpcParser{
		schema:          i.Schema,
		goTypeParser:    i.GoTypeParser,
		templateManager: templateManager,
	}
}
