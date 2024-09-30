package grpc

import (
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

	GoTypeParser    types_parser.TypeParser
	TemplateManager template.TemplateManager
}

func NewGrpcParser(i *NewGrpcParserInput) GrpcParser {
	return &goGrpcParser{
		schema:          i.Schema,
		goTypeParser:    i.GoTypeParser,
		templateManager: i.TemplateManager,
	}
}
