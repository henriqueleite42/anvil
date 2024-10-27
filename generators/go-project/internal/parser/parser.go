package parser

import (
	generator_config "github.com/henriqueleite42/anvil/generators/go-project/config"
	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/formatter"
	"github.com/henriqueleite42/anvil/language-helpers/golang/imports"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

type ParserRepository struct {
	Methods []*templates.TemplMethod
}

type ParserUsecase struct {
	Methods []*templates.TemplMethod
}

type ParserGrpcDelivery struct {
	Methods []*templates.TemplMethodDelivery
}

type Parser struct {
	schema *schemas.AnvpSchema
	config *generator_config.GeneratorConfig

	goTypesParser types_parser.TypesParser

	repositories   map[string]*ParserRepository
	usecases       map[string]*ParserUsecase
	grpcDeliveries map[string]*ParserGrpcDelivery
}

func NewTypesParser(
	schema *schemas.AnvpSchema,
	config *generator_config.GeneratorConfig,
) (*Parser, error) {
	modelsImport := imports.NewImport(config.ModuleName+"/internal/models", nil)

	goTypesParser, err := types_parser.NewTypeParser(&types_parser.NewTypeParserInput{
		Schema: schema,
		GetEnumsImport: func(e *schemas.Enum) *imports.Import {
			return modelsImport
		},
		GetTypesImport: func(t *schemas.Type) *imports.Import {
			return modelsImport
		},
		GetEventsImport: func(t *schemas.Type) *imports.Import {
			return modelsImport
		},
		GetEntitiesImport: func(t *schemas.Type) *imports.Import {
			return modelsImport
		},
		GetRepositoryImport: func(t *schemas.Type) *imports.Import {
			domainSnake := formatter.PascalToSnake(t.Domain)
			path := config.ModuleName + "/internal/repository/" + domainSnake
			alias := domainSnake + "_repository"

			return imports.NewImport(path, &alias)
		},
		GetUsecaseImport: func(t *schemas.Type) *imports.Import {
			domainSnake := formatter.PascalToSnake(t.Domain)
			path := config.ModuleName + "/internal/usecase/" + domainSnake
			alias := domainSnake + "_usecase"

			return imports.NewImport(path, &alias)
		},
	})
	if err != nil {
		return nil, err
	}

	return &Parser{
		schema: schema,
		config: config,

		goTypesParser: goTypesParser,

		repositories:   map[string]*ParserRepository{},
		usecases:       map[string]*ParserUsecase{},
		grpcDeliveries: map[string]*ParserGrpcDelivery{},
	}, nil
}
