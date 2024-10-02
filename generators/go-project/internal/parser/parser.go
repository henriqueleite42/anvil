package parser

import (
	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

type Parser struct {
	ModelsPath string

	Schema *schemas.Schema

	GoTypesParserModels     types_parser.TypeParser
	GoTypesParserRepository types_parser.TypeParser
	GoTypesParserUsecase    types_parser.TypeParser

	MethodsRepository                []*templates.TemplMethod
	MethodsUsecaseToAvoidDuplication map[string]bool
	MethodsUsecase                   []*templates.TemplMethod
	MethodsGrpcDelivery              []*templates.TemplMethodDelivery
}
