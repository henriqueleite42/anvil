package parser

import (
	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

type Parser struct {
	Schema *schemas.AnvpSchema

	GoTypesParserModels     types_parser.TypeParser
	GoTypesParserRepository types_parser.TypeParser
	GoTypesParserUsecase    types_parser.TypeParser

	MethodsRepository                []*templates.TemplMethod
	MethodsUsecaseToAvoidDuplication map[string]bool
	MethodsUsecase                   []*templates.TemplMethod
	MethodsGrpcDelivery              []*templates.TemplMethodDelivery
}

func hasModels(schema *schemas.AnvpSchema, curDomain string) bool {
	if schema.Types != nil && schema.Types.Types != nil {
		if _, ok := schema.Types.Types[curDomain]; ok {
			return true
		}
	}
	if schema.Entities != nil && schema.Entities.Entities != nil {
		return true
	}
	if schema.Events != nil && schema.Events.Events != nil {
		return true
	}
	if schema.Enums != nil && schema.Enums.Enums != nil {
		return true
	}

	return false
}

func hasRepositories(schema *schemas.AnvpSchema, curDomain string) bool {
	if schema.Repositories != nil && schema.Repositories.Repositories != nil {
		if _, ok := schema.Repositories.Repositories[curDomain]; ok {
			return true
		}
	}

	return false
}

func hasUsecases(schema *schemas.AnvpSchema, curDomain string) bool {
	if schema.Usecases != nil && schema.Usecases.Usecases != nil {
		if _, ok := schema.Usecases.Usecases[curDomain]; ok {
			return true
		}
	}

	return false
}

func NewTypesParser(
	schema *schemas.AnvpSchema,
	curDomain string,
	parserInput *types_parser.NewTypeParserInput,
) (*Parser, error) {
	typeParser := &Parser{
		Schema: schema,
	}

	if hasModels(schema, curDomain) {
		goTypesParserModels, err := types_parser.NewTypeParser(parserInput)
		if err != nil {
			return nil, err
		}

		typeParser.GoTypesParserModels = goTypesParserModels
	}

	if hasRepositories(schema, curDomain) {
		goTypesParserRepository, err := types_parser.NewTypeParser(parserInput)
		if err != nil {
			return nil, err
		}

		typeParser.GoTypesParserRepository = goTypesParserRepository
		typeParser.MethodsRepository = make(
			[]*templates.TemplMethod,
			0,
			len(schema.Repositories.Repositories[curDomain].Methods.Methods),
		)

		goTypesParserRepository.AddImport("context")
	}

	if hasUsecases(schema, curDomain) {
		goTypesParserUsecase, err := types_parser.NewTypeParser(parserInput)
		if err != nil {
			return nil, err
		}

		typeParser.GoTypesParserUsecase = goTypesParserUsecase
		typeParser.MethodsUsecase = make(
			[]*templates.TemplMethod,
			0,
			len(schema.Usecases.Usecases[curDomain].Methods.Methods),
		)
	}

	return typeParser, nil
}
