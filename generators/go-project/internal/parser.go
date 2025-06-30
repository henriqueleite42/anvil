package internal

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"
	generator_config "github.com/henriqueleite42/anvil/generators/go-project/config"
	"github.com/henriqueleite42/anvil/generators/go-project/internal/parser"
	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/imports"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	"github.com/henriqueleite42/anvil/language-helpers/golang/template"
)

var templatesNamesValues = map[string]string{
	"models":                      templates.ModelsTempl,
	"repository":                  templates.RepositoryTempl,
	"repository-struct":           templates.RepositoryStructTempl,
	"repository-method":           templates.RepositoryMethodTempl,
	"usecase":                     templates.UsecaseTempl,
	"usecase-struct":              templates.UsecaseStructTempl,
	"usecase-method":              templates.UsecaseMethodTempl,
	"grpc-delivery-module-helper": templates.GrpcDeliveryModuleHelperTempl,
	"grpc-delivery-module":        templates.GrpcDeliveryModuleTempl,
	"http-delivery-module":        templates.HttpDeliveryModuleTempl,
	"http-delivery-route":         templates.HttpDeliveryRouteTempl,
	"http-delivery":               templates.HttpDeliveryTempl,
	"queue-delivery-module":       templates.QueueDeliveryModuleTempl,
	"queue-delivery":              templates.QueueDeliveryTempl,
	"go-mod":                      templates.GoModTempl,
	"validator-implementation":    templates.ValidatorImplementationTempl,
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

	typeParser, err := parser.NewTypesParser(schema, config)
	if err != nil {
		return nil, err
	}

	err = typeParser.Parse()
	if err != nil {
		return nil, err
	}

	enumsPerDomain, err := typeParser.GetEnums()
	if err != nil {
		return nil, err
	}
	typesPerDomain, err := typeParser.GetTypes()
	if err != nil {
		return nil, err
	}
	entitiesPerDomain, err := typeParser.GetEntities()
	if err != nil {
		return nil, err
	}
	eventsPerDomain, err := typeParser.GetEvents()
	if err != nil {
		return nil, err
	}
	repositoryTypesPerDomain, err := typeParser.GetRepositoryTypes()
	if err != nil {
		return nil, err
	}
	repositoriesPerDomain, err := typeParser.GetRepositories()
	if err != nil {
		return nil, err
	}
	usecaseTypesPerDomain, err := typeParser.GetUsecaseTypes()
	if err != nil {
		return nil, err
	}
	usecasesPerDomain, err := typeParser.GetUsecases()
	if err != nil {
		return nil, err
	}
	grpcDeliveriesPerDomain, err := typeParser.GetGrpcDeliveries()
	if err != nil {
		return nil, err
	}
	httpDeliveriesPerDomain, err := typeParser.GetHttpDeliveries()
	if err != nil {
		return nil, err
	}
	queueDeliveriesPerDomain, err := typeParser.GetQueueDeliveries()
	if err != nil {
		return nil, err
	}

	var hasDelivery bool
	var hasGrpcDelivery bool
	var hasHttpDelivery bool
	var hasQueueDelivery bool

	for _, scm := range schema.Schemas {
		var hasModels bool
		var hasRepository bool
		var hasUsecase bool
		var domainHasEnums bool
		var domainHasGrpcDelivery bool
		var domainHasHttpDelivery bool
		var domainHasQueueDelivery bool

		curDomain := scm.Domain

		templInput := &templates.TemplInput{
			ProjectName:                 config.ProjectName,
			DomainPascal:                curDomain,
			DomainSnake:                 strcase.ToSnake(curDomain),
			DomainCamel:                 strcase.ToCamel(curDomain),
			SpacingRelativeToDomainName: strings.Repeat(" ", len(curDomain)-1),
		}

		// Models

		if enumsPerDomain != nil {
			if enums, ok := enumsPerDomain[curDomain]; ok {
				hasModels = true
				domainHasEnums = true
				templInput.Enums = enums
			}
		}
		if typesPerDomain != nil {
			if types, ok := typesPerDomain[curDomain]; ok {
				for _, t := range types {
					typeParser.ImportsModels[curDomain].MergeImports(t.ImportsUnorganized)
				}

				hasModels = true
				templInput.Types = types
			}
		}
		if eventsPerDomain != nil {
			if events, ok := eventsPerDomain[curDomain]; ok {
				for _, t := range events {
					typeParser.ImportsModels[curDomain].MergeImports(t.ImportsUnorganized)
				}

				hasModels = true
				templInput.Events = events
			}
		}
		if entitiesPerDomain != nil {
			if entities, ok := entitiesPerDomain[curDomain]; ok {
				for _, t := range entities {
					typeParser.ImportsModels[curDomain].MergeImports(t.ImportsUnorganized)
				}

				hasModels = true
				templInput.Entities = entities
			}
		}

		importsModels := imports.ResolveImports(
			typeParser.ImportsModels[curDomain].GetImportsUnorganized(),
			"models",
		)
		templInput.ImportsModels = importsModels

		// Repository

		if repositoryTypesPerDomain != nil {
			if repositoryTypes, ok := repositoryTypesPerDomain[curDomain]; ok {
				hasRepository = true
				for _, r := range repositoryTypes {
					typeParser.ImportsRepository[curDomain].MergeImports(r.ImportsUnorganized)
				}

				templInput.TypesRepository = repositoryTypes
			}
		}
		if repositoriesPerDomain != nil {
			if repository, ok := repositoriesPerDomain[curDomain]; ok {
				hasRepository = true
				templInput.MethodsRepository = repository.Methods
			}
		}

		importsRepository := imports.ResolveImports(
			typeParser.ImportsRepository[curDomain].GetImportsUnorganized(),
			templInput.DomainSnake+"_repository",
		)
		templInput.ImportsRepository = importsRepository

		// Usecase

		if usecaseTypesPerDomain != nil {
			if usecaseTypes, ok := usecaseTypesPerDomain[curDomain]; ok {
				for _, r := range usecaseTypes {
					typeParser.ImportsUsecase[curDomain].MergeImports(r.ImportsUnorganized)
				}

				hasUsecase = true
				templInput.TypesUsecase = usecaseTypes
			}
		}
		if usecasesPerDomain != nil {
			if usecase, ok := usecasesPerDomain[curDomain]; ok {
				hasUsecase = true
				templInput.MethodsUsecase = usecase.Methods
			}
		}

		importsUsecase := imports.ResolveImports(
			typeParser.ImportsUsecase[curDomain].GetImportsUnorganized(),
			templInput.DomainSnake+"_usecase",
		)
		templInput.ImportsUsecase = importsUsecase

		// Grpc Delivery
		typeParser.ImportsGrpcDeliveryHelper[curDomain].MergeImport(config.PbModuleImport)

		if grpcDeliveriesPerDomain != nil {
			if grpcDelivery, ok := grpcDeliveriesPerDomain[curDomain]; ok {
				hasDelivery = true
				hasGrpcDelivery = true
				domainHasGrpcDelivery = true

				for _, method := range grpcDelivery.Methods {
					if method.Input != nil {
						typeParser.ImportsGrpcDelivery[curDomain].MergeImports(method.Input.ImportsUnorganized)
					} else {
						typeParser.ImportsGrpcDelivery[curDomain].AddImport("google.golang.org/protobuf/types/known/emptypb", nil)
					}
					if method.Output != nil {
						typeParser.ImportsGrpcDelivery[curDomain].MergeImports(method.Output.ImportsUnorganized)
					} else {
						typeParser.ImportsGrpcDelivery[curDomain].AddImport("google.golang.org/protobuf/types/known/emptypb", nil)
					}
				}

				typeParser.ImportsGrpcDelivery[curDomain].AddImport("context", nil)
				typeParser.ImportsGrpcDelivery[curDomain].AddImport("github.com/rs/xid", nil)
				typeParser.ImportsGrpcDelivery[curDomain].AddImport("google.golang.org/grpc", nil)
				typeParser.ImportsGrpcDelivery[curDomain].AddImport("github.com/rs/zerolog", nil)
				typeParser.ImportsGrpcDelivery[curDomain].MergeImport(config.PbModuleImport)
				typeParser.ImportsGrpcDelivery[curDomain].AddImport(config.ProjectName+"/internal/adapters", nil)

				if _, ok := enumsPerDomain[curDomain]; ok {
					typeParser.ImportsGrpcDeliveryHelper[curDomain].AddImport(config.ProjectName+"/internal/models", nil)
				}

				templInput.MethodsGrpcDelivery = grpcDelivery.Methods
			}
		}

		importsGrpcDelivery := imports.ResolveImports(
			typeParser.ImportsGrpcDelivery[curDomain].GetImportsUnorganized(),
			templInput.DomainSnake+"_grpc_delivery",
		)
		templInput.ImportsGrpcDelivery = importsGrpcDelivery
		importsGrpcDeliveryHelper := imports.ResolveImports(
			typeParser.ImportsGrpcDeliveryHelper[curDomain].GetImportsUnorganized(),
			templInput.DomainSnake+"_grpc_delivery_helper",
		)
		templInput.ImportsGrpcDeliveryHelper = importsGrpcDeliveryHelper

		// Http Delivery
		if httpDeliveriesPerDomain != nil {
			if httpDelivery, ok := httpDeliveriesPerDomain[curDomain]; ok {
				hasDelivery = true
				hasHttpDelivery = true
				domainHasHttpDelivery = true

				typeParser.ImportsHttpDelivery[curDomain].AddImport("net/http", nil)
				typeParser.ImportsHttpDelivery[curDomain].AddImport("github.com/rs/zerolog", nil)
				typeParser.ImportsHttpDelivery[curDomain].AddImport(config.ProjectName+"/internal/adapters", nil)
				usecaseAlias := templInput.DomainSnake + "_usecase"
				typeParser.ImportsHttpDelivery[curDomain].AddImport(config.ProjectName+"/internal/usecase/"+templInput.DomainSnake, &usecaseAlias)

				templInput.MethodsHttpDelivery = httpDelivery.Methods
			}
		}

		importsHttpDelivery := imports.ResolveImports(
			typeParser.ImportsHttpDelivery[curDomain].GetImportsUnorganized(),
			templInput.DomainSnake+"_http_delivery",
		)
		templInput.ImportsHttpDelivery = importsHttpDelivery

		// Queue Delivery
		if queueDeliveriesPerDomain != nil {
			if queueDelivery, ok := queueDeliveriesPerDomain[curDomain]; ok {
				hasDelivery = true
				hasQueueDelivery = true
				domainHasQueueDelivery = true

				for _, method := range queueDelivery.Methods {
					if method.Input != nil {
						typeParser.ImportsQueueDelivery[curDomain].MergeImport(method.Input.ModuleImport)
					}
				}

				typeParser.ImportsQueueDelivery[curDomain].AddImport("context", nil)
				typeParser.ImportsQueueDelivery[curDomain].AddImport("encoding/json", nil)
				typeParser.ImportsQueueDelivery[curDomain].AddImport(config.ProjectName+"/internal/adapters", nil)

				templInput.MethodsQueueDelivery = queueDelivery.Methods
			}
		}

		importsQueueDelivery := imports.ResolveImports(
			typeParser.ImportsQueueDelivery[curDomain].GetImportsUnorganized(),
			templInput.DomainSnake+"_queue_delivery",
		)
		templInput.ImportsQueueDelivery = importsQueueDelivery

		// -----------
		//
		// Resolve files
		//
		// -----------

		if hasModels {
			models, err := templateManager.Parse("models", templInput)
			if err != nil {
				return nil, err
			}
			files = append(files, &File{
				Name:      "internal/models/" + templInput.DomainSnake + ".go",
				Content:   models,
				Overwrite: true,
			})
		}

		if hasRepository {
			repositoryContract, err := templateManager.Parse("repository", templInput)
			if err != nil {
				return nil, err
			}
			files = append(files, &File{
				Name:      fmt.Sprintf("internal/repository/%s/contract.go", templInput.DomainSnake),
				Content:   repositoryContract,
				Overwrite: true,
			})

			repositoryStruct, err := templateManager.Parse("repository-struct", templInput)
			if err != nil {
				return nil, err
			}
			files = append(files, &File{
				Name:    fmt.Sprintf("internal/repository/%s/implementation.go", templInput.DomainSnake),
				Content: repositoryStruct,
			})

			for _, v := range templInput.MethodsRepository {
				repositoryMethod, err := templateManager.Parse("repository-method", &templates.RepositoryMethodTemplInput{
					DomainPascal:   templInput.DomainPascal,
					DomainCamel:    templInput.DomainCamel,
					DomainSnake:    templInput.DomainSnake,
					MethodName:     v.MethodName,
					InputTypeName:  v.InputTypeName,
					OutputTypeName: v.OutputTypeName,
					Imports:        v.Imports,
				})
				if err != nil {
					return nil, err
				}
				methodNameSnake := strcase.ToSnake(v.MethodName)
				files = append(files, &File{
					Name:    fmt.Sprintf("internal/repository/%s/%s.go", templInput.DomainSnake, methodNameSnake),
					Content: repositoryMethod,
				})
			}
		}

		if hasUsecase {
			usecaseContract, err := templateManager.Parse("usecase", templInput)
			if err != nil {
				return nil, err
			}
			files = append(files, &File{
				Name:      fmt.Sprintf("internal/usecase/%s/contract.go", templInput.DomainSnake),
				Content:   usecaseContract,
				Overwrite: true,
			})

			usecaseStruct, err := templateManager.Parse("usecase-struct", templInput)
			if err != nil {
				return nil, err
			}
			files = append(files, &File{
				Name:    fmt.Sprintf("internal/usecase/%s/implementation.go", templInput.DomainSnake),
				Content: usecaseStruct,
			})

			for _, v := range templInput.MethodsUsecase {
				usecaseMethod, err := templateManager.Parse("usecase-method", &templates.UsecaseMethodTemplInput{
					DomainPascal:   templInput.DomainPascal,
					DomainCamel:    templInput.DomainCamel,
					DomainSnake:    templInput.DomainSnake,
					MethodName:     v.MethodName,
					InputTypeName:  v.InputTypeName,
					OutputTypeName: v.OutputTypeName,
					Imports:        v.Imports,
				})
				if err != nil {
					return nil, err
				}
				methodNameSnake := strcase.ToSnake(v.MethodName)
				files = append(files, &File{
					Name:    fmt.Sprintf("internal/usecase/%s/%s.go", templInput.DomainSnake, methodNameSnake),
					Content: usecaseMethod,
				})
			}
		}

		if domainHasGrpcDelivery {
			grpcModule, err := templateManager.Parse("grpc-delivery-module", templInput)
			if err != nil {
				return nil, err
			}
			files = append(
				files,
				&File{
					Name:      fmt.Sprintf("internal/delivery/grpc/%s/%s.go", templInput.DomainSnake, templInput.DomainSnake),
					Content:   grpcModule,
					Overwrite: true,
				},
			)

			if domainHasEnums {
				grpcModuleHelper, err := templateManager.Parse("grpc-delivery-module-helper", templInput)
				if err != nil {
					return nil, err
				}
				files = append(
					files,
					&File{
						Name:      fmt.Sprintf("internal/delivery/grpc/%s/helpers/enums.go", templInput.DomainSnake),
						Content:   grpcModuleHelper,
						Overwrite: true,
					},
				)
			}
		}

		if domainHasHttpDelivery {
			httpModule, err := templateManager.Parse("http-delivery-module", templInput)
			if err != nil {
				return nil, err
			}
			files = append(
				files,
				&File{
					Name:    fmt.Sprintf("internal/delivery/http/%s/%s_controller.go", templInput.DomainSnake, templInput.DomainSnake),
					Content: httpModule,
				},
			)

			for _, route := range templInput.MethodsHttpDelivery {
				httpRoute, err := templateManager.Parse("http-delivery-route", route)
				if err != nil {
					return nil, err
				}
				files = append(
					files,
					&File{
						Name:    fmt.Sprintf("internal/delivery/http/%s/%s.go", route.DomainSnake, route.RouteNameSnake),
						Content: httpRoute,
					},
				)
			}
		}

		if domainHasQueueDelivery {
			queueModule, err := templateManager.Parse("queue-delivery-module", templInput)
			if err != nil {
				return nil, err
			}
			files = append(
				files,
				&File{
					Name:    fmt.Sprintf("internal/delivery/queue/%s.go", templInput.DomainSnake),
					Content: queueModule,
				},
			)
		}
	}

	if hasDelivery {
		validatorImplementation, err := templateManager.Parse("validator-implementation", config)
		if err != nil {
			return nil, err
		}

		files = append(
			files,
			&File{
				Name:    "internal/adapters/validator.go",
				Content: templates.ValidatorTempl,
			},
			&File{
				Name:    "internal/adapters/go-validator/go-validator.go",
				Content: validatorImplementation,
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

	if hasHttpDelivery {
		httpImplementation, err := templateManager.Parse("http-delivery", config)
		if err != nil {
			return nil, err
		}

		files = append(
			files,
			&File{
				Name:    "internal/delivery/http/http.go",
				Content: httpImplementation,
			},
		)
	}

	if hasQueueDelivery {
		queueImplementation, err := templateManager.Parse("queue-delivery", config)
		if err != nil {
			return nil, err
		}

		files = append(
			files,
			&File{
				Name:    "internal/delivery/queue/queue.go",
				Content: queueImplementation,
			},
		)
	}

	goConfig, err := templateManager.Parse("go-mod", &templates.GoConfig{
		PkgName:   config.ProjectName,
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
