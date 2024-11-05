package grpc

import (
	"github.com/henriqueleite42/anvil/language-helpers/golang/grpc/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/imports"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	"github.com/henriqueleite42/anvil/language-helpers/golang/template"
	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

type goGrpcParser struct {
	schema *schemas.AnvpSchema

	goTypeParser    types_parser.TypesParser
	templateManager template.TemplateManager

	getEnumConversionImpt func(e *schemas.Enum) *imports.Import
}

type NewGrpcParserInput struct {
	Schema *schemas.AnvpSchema

	GoTypeParser types_parser.TypesParser

	GetEnumConversionImpt func(e *schemas.Enum) *imports.Import
}

func NewGrpcParser(i *NewGrpcParserInput) GrpcParser {
	templateManager := template.NewTemplateManager()

	templateManager.AddTemplate("input-prop-list", templates.InputPropListTempl)
	templateManager.AddTemplate("input-prop-map", templates.InputPropMapTempl)
	templateManager.AddTemplate("input-prop-optional", templates.InputPropOptionalTempl)

	return &goGrpcParser{
		schema: i.Schema,

		goTypeParser:    i.GoTypeParser,
		templateManager: templateManager,

		getEnumConversionImpt: i.GetEnumConversionImpt,
	}
}
