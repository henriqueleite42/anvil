package internal

import (
	"fmt"
	"strings"

	generator_config "github.com/henriqueleite42/anvil/generators/go-project/config"
	"github.com/henriqueleite42/anvil/generators/go-project/internal/parser"
	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/formatter"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	"github.com/henriqueleite42/anvil/language-helpers/golang/template"
	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

var templatesNamesValues = map[string]string{
	"models":               templates.ModelsTempl,
	"repository":           templates.RepositoryTempl,
	"repository-struct":    templates.RepositoryStructTempl,
	"repository-method":    templates.RepositoryMethodTempl,
	"usecase":              templates.UsecaseTempl,
	"usecase-struct":       templates.UsecaseStructTempl,
	"usecase-method":       templates.UsecaseMethodTempl,
	"grpc-delivery-module": templates.GrpcDeliveryModuleTempl,
	"go-mod":               templates.GoModTempl,
}

type File struct {
	Name      string
	Content   string
	Overwrite bool
}

func Parse(schema *schemas.AnvpSchema, config *generator_config.GeneratorConfig) ([]*File, error) {
	files := []*File{}

	templateManager := template.NewTemplateManager()
	for k, v := range templatesNamesValues {
		err := templateManager.AddTemplate(k, v)
		if err != nil {
			return nil, err
		}
	}

	var hasDelivery bool
	var hasGrpcDelivery bool
	for _, v := range schema.Schemas {
		curDomain := v.Domain
		domainSnake := formatter.PascalToSnake(curDomain)

		typeParser, err := parser.NewTypesParser(schema, curDomain, &types_parser.NewTypeParserInput{
			Schema:        schema,
			EnumsPkg:      "models",
			TypesPkg:      "models",
			EventsPkg:     "models",
			EntitiesPkg:   "models",
			RepositoryPkg: domainSnake + "_repository",
			UsecasePkg:    domainSnake + "_usecase",
		})
		if err != nil {
			return nil, err
		}

		// Parse

		models, err := typeParser.ParseModels(curDomain)
		if err != nil {
			return nil, err
		}

		repositories, err := typeParser.ParseRepositories(curDomain)
		if err != nil {
			return nil, err
		}

		usecases, err := typeParser.ParseUsecases(curDomain)
		if err != nil {
			return nil, err
		}

		deliveryGrpc, err := typeParser.ParseDeliveriesGrpc(curDomain)
		if err != nil {
			return nil, err
		}

		// Templates

		templInput := &templates.TemplInput{
			Domain:                      curDomain,
			DomainSnake:                 domainSnake,
			DomainCamel:                 formatter.PascalToCamel(curDomain),
			SpacingRelativeToDomainName: strings.Repeat(" ", len(curDomain)),
			ImportsModels:               models.Imports,
			ImportsRepository:           repositories.Imports,
			ImportsUsecase:              usecases.Imports,
			ImportsGrpcDelivery:         deliveryGrpc.Imports,
			Enums:                       models.Enums,
			Entities:                    models.Entities,
			Events:                      models.Events,
			TypesRepository:             repositories.Types,
			TypesUsecase:                usecases.Types,
			MethodsRepository:           typeParser.MethodsRepository,
			MethodsUsecase:              typeParser.MethodsUsecase,
			MethodsGrpcDelivery:         typeParser.MethodsGrpcDelivery,
		}

		domainKebab := formatter.PascalToKebab(curDomain)

		if templInput.Enums != nil ||
			templInput.Entities != nil ||
			templInput.Events != nil {
			models, err := templateManager.Parse("models", templInput)
			if err != nil {
				return nil, err
			}
			files = append(files, &File{
				Name:      "internal/models/" + domainKebab + ".go",
				Content:   models,
				Overwrite: true,
			})
		}

		if templInput.MethodsRepository != nil {
			repository, err := templateManager.Parse("repository", templInput)
			if err != nil {
				return nil, err
			}
			files = append(files, &File{
				Name:      fmt.Sprintf("internal/repository/%s/%s.go", domainKebab, domainKebab),
				Content:   repository,
				Overwrite: true,
			})

			repositoryStruct, err := templateManager.Parse("repository-struct", templInput)
			if err != nil {
				return nil, err
			}
			files = append(files, &File{
				Name:    fmt.Sprintf("internal/repository/%s/implementation.go", domainKebab),
				Content: repositoryStruct,
			})

			for _, v := range templInput.MethodsRepository {
				repositoryMethod, err := templateManager.Parse("repository-method", &templates.RepositoryMethodTemplInput{
					Domain:         templInput.Domain,
					DomainSnake:    templInput.DomainSnake,
					MethodName:     v.MethodName,
					InputTypeName:  v.InputTypeName,
					OutputTypeName: v.OutputTypeName,
				})
				if err != nil {
					return nil, err
				}
				methodNameKebab := formatter.PascalToKebab(v.MethodName)
				files = append(files, &File{
					Name:    fmt.Sprintf("internal/repository/%s/%s.go", domainKebab, methodNameKebab),
					Content: repositoryMethod,
				})
			}
		}

		if templInput.MethodsUsecase != nil {
			usecase, err := templateManager.Parse("usecase", templInput)
			if err != nil {
				return nil, err
			}
			files = append(files, &File{
				Name:      fmt.Sprintf("internal/usecase/%s/%s.go", domainKebab, domainKebab),
				Content:   usecase,
				Overwrite: true,
			})

			usecaseStruct, err := templateManager.Parse("usecase-struct", templInput)
			if err != nil {
				return nil, err
			}
			files = append(files, &File{
				Name:    fmt.Sprintf("internal/usecase/%s/implementation.go", domainKebab),
				Content: usecaseStruct,
			})

			for _, v := range templInput.MethodsUsecase {
				usecaseMethod, err := templateManager.Parse("usecase-method", &templates.UsecaseMethodTemplInput{
					Domain:         templInput.Domain,
					DomainSnake:    templInput.DomainSnake,
					MethodName:     v.MethodName,
					InputTypeName:  v.InputTypeName,
					OutputTypeName: v.OutputTypeName,
				})
				if err != nil {
					return nil, err
				}
				methodNameKebab := formatter.PascalToKebab(v.MethodName)
				files = append(files, &File{
					Name:    fmt.Sprintf("internal/usecase/%s/%s.go", domainKebab, methodNameKebab),
					Content: usecaseMethod,
				})
			}
		}

		if templInput.MethodsGrpcDelivery != nil {
			hasDelivery = true
			hasGrpcDelivery = true

			grpcModule, err := templateManager.Parse("grpc-delivery-module", templInput)
			if err != nil {
				return nil, err
			}
			files = append(
				files,
				&File{
					Name:      fmt.Sprintf("internal/delivery/grpc/%s/%s.go", templInput.DomainSnake, domainKebab),
					Content:   grpcModule,
					Overwrite: true,
				},
			)
		}
	}

	if hasDelivery {
		files = append(
			files,
			&File{
				Name:      "internal/adapters/validator.go",
				Content:   templates.ValidatorTempl,
				Overwrite: true,
			},
			&File{
				Name:    "internal/adapters/go-validator/go-validator.go",
				Content: templates.ValidatorImplementationTempl,
			},
			&File{
				Name:      "internal/delivery/delivery.go",
				Content:   templates.DeliveryTempl,
				Overwrite: true,
			},
			&File{
				Name:      "internal/delivery/gratefully-shutdown.go",
				Content:   templates.GratefullyShutdownTempl,
				Overwrite: true,
			},
			&File{
				Name:      "internal/utils/sync.go",
				Content:   templates.UtilsSyncTempl,
				Overwrite: true,
			},
		)
	}

	if hasGrpcDelivery {
		files = append(
			files,
			&File{
				Name:    "internal/delivery/grpc/grpc.go",
				Content: templates.GrpcDeliveryTempl,
			},
		)
	}

	goConfig, err := templateManager.Parse("go-mod", &templates.GoConfig{
		PkgName:   config.ModuleName,
		GoVersion: config.GoVersion,
	})
	if err != nil {
		return nil, err
	}

	files = append(
		files,
		&File{
			Name:    "cmd/main.go",
			Content: templates.MainTempl,
		},
		&File{
			Name:    "go.mod",
			Content: goConfig,
		},
		&File{
			Name:    ".editorconfig",
			Content: templates.EditorConfigTempl,
		},
		&File{
			Name:    ".gitignore",
			Content: templates.GitIgnoreTempl,
		},
		&File{
			Name:    "README.md",
			Content: templates.ReadMeTempl,
		},
		&File{
			Name:    "staticcheck.conf",
			Content: templates.StaticCheckTempl,
		},
	)

	return files, nil
}
