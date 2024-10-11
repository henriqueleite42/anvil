package internal

import (
	"fmt"
	"sort"
	"strings"

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
}

type File struct {
	Name      string
	Content   string
	Overwrite bool
}

func Parse(schema *schemas.Schema) ([]*File, error) {
	lenEnums := 0
	if schema.Enums != nil && schema.Enums.Enums != nil {
		lenEnums = len(schema.Enums.Enums)
	}
	lenEntities := 0
	if schema.Entities != nil && schema.Entities.Entities != nil {
		lenEntities = len(schema.Entities.Entities)
	}
	lenRepository := 0
	if schema.Repository != nil && schema.Repository.Methods != nil && schema.Repository.Methods.Methods != nil {
		lenRepository = len(schema.Repository.Methods.Methods)
	}
	lenUsecase := 0
	if schema.Usecase != nil && schema.Usecase.Methods != nil && schema.Usecase.Methods.Methods != nil {
		lenUsecase = len(schema.Usecase.Methods.Methods)
	}
	lenGrpcDelivery := 0
	if schema.Delivery != nil &&
		schema.Delivery.Grpc != nil &&
		schema.Delivery.Grpc.Rpcs != nil &&
		schema.Usecase != nil &&
		schema.Usecase.Methods != nil &&
		schema.Usecase.Methods.Methods != nil {
		lenGrpcDelivery = len(schema.Delivery.Grpc.Rpcs)
	}

	domainSnake := formatter.PascalToSnake(schema.Domain)

	goTypesParserModels, err := types_parser.NewTypeParser(&types_parser.NewTypeParserInput{
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
	goTypesParserRepository, err := types_parser.NewTypeParser(&types_parser.NewTypeParserInput{
		Schema:        schema,
		EnumsPkg:      "models",
		TypesPkg:      domainSnake + "_repository",
		EventsPkg:     "models",
		EntitiesPkg:   "models",
		RepositoryPkg: domainSnake + "_repository",
		UsecasePkg:    domainSnake + "_usecase",
	})
	if err != nil {
		return nil, err
	}
	goTypesParserUsecase, err := types_parser.NewTypeParser(&types_parser.NewTypeParserInput{
		Schema:        schema,
		EnumsPkg:      "models",
		TypesPkg:      domainSnake + "_usecase",
		EventsPkg:     "models",
		EntitiesPkg:   "models",
		RepositoryPkg: domainSnake + "_repository",
		UsecasePkg:    domainSnake + "_usecase",
	})
	if err != nil {
		return nil, err
	}

	goTypesParserRepository.AddImport("context")

	typeParser := &parser.Parser{
		ModelsPath: "foo/models",
		Schema:     schema,

		GoTypesParserModels:     goTypesParserModels,
		GoTypesParserRepository: goTypesParserRepository,
		GoTypesParserUsecase:    goTypesParserUsecase,

		MethodsUsecaseToAvoidDuplication: map[string]bool{},

		MethodsRepository:   make([]*templates.TemplMethod, 0, lenRepository),
		MethodsUsecase:      make([]*templates.TemplMethod, 0, lenUsecase),
		MethodsGrpcDelivery: make([]*templates.TemplMethodDelivery, 0, lenGrpcDelivery),
	}

	if lenEnums > 0 {
		for _, v := range schema.Enums.Enums {
			_, err := goTypesParserModels.ParseEnum(v)
			if err != nil {
				return nil, err
			}
		}
	}

	if lenEntities > 0 {
		for _, v := range schema.Entities.Entities {
			err := typeParser.ResolveEntity(v)
			if err != nil {
				return nil, err
			}
		}
	}

	if lenRepository > 0 {
		for _, v := range schema.Repository.Methods.Methods {
			err := typeParser.ResolveRepositoryMethod(v, domainSnake+"_repository")
			if err != nil {
				return nil, err
			}
		}
	}

	// As the usecase is used by all the deliveries,
	// we reuse the type parsing for all of them,
	// but reset the imports to be sure to import the right things
	var importsUsecase [][]string = nil
	if lenUsecase > 0 {
		for _, v := range schema.Usecase.Methods.Methods {
			err := typeParser.ResolveUsecaseMethod(v, domainSnake+"_usecase")
			if err != nil {
				return nil, err
			}
		}

		goTypesParserUsecase.AddImport("context")
		importsUsecase = goTypesParserUsecase.GetImports()
		goTypesParserUsecase.ResetImports()
	}

	var importsGrpcDelivery [][]string = nil
	if lenGrpcDelivery > 0 {
		for _, v := range schema.Delivery.Grpc.Rpcs {
			err := typeParser.ResolveGrpcDelivery(v)
			if err != nil {
				return nil, err
			}
		}

		goTypesParserUsecase.AddImport("context")
		goTypesParserUsecase.AddImport("github.com/rs/xid")
		goTypesParserUsecase.AddImport("github.com/rs/zerolog")
		goTypesParserUsecase.AddImport("google.golang.org/grpc")
		importsGrpcDelivery = goTypesParserUsecase.GetImports()
		goTypesParserUsecase.ResetImports()
	}

	importsModels := goTypesParserModels.GetImports()
	importsRepository := goTypesParserRepository.GetImports()

	enums := goTypesParserModels.GetEnums()
	entities := goTypesParserModels.GetEntities()
	typesRepository := goTypesParserRepository.GetRepository()
	typesRepository = append(typesRepository, goTypesParserRepository.GetTypes()...)
	sort.Slice(typesRepository, func(i, j int) bool {
		return typesRepository[i].GolangType < typesRepository[j].GolangType
	})
	typesUsecase := goTypesParserUsecase.GetUsecase()
	typesUsecase = append(typesUsecase, goTypesParserUsecase.GetTypes()...)
	sort.Slice(typesUsecase, func(i, j int) bool {
		return typesUsecase[i].GolangType < typesUsecase[j].GolangType
	})

	sort.Slice(typeParser.MethodsRepository, func(i, j int) bool {
		return typeParser.MethodsRepository[i].Order < typeParser.MethodsRepository[j].Order
	})
	sort.Slice(typeParser.MethodsUsecase, func(i, j int) bool {
		return typeParser.MethodsUsecase[i].Order < typeParser.MethodsUsecase[j].Order
	})
	sort.Slice(typeParser.MethodsGrpcDelivery, func(i, j int) bool {
		return typeParser.MethodsGrpcDelivery[i].Order < typeParser.MethodsGrpcDelivery[j].Order
	})

	templInput := &templates.TemplInput{
		Domain:                      schema.Domain,
		DomainSnake:                 domainSnake,
		DomainCamel:                 formatter.PascalToCamel(schema.Domain),
		SpacingRelativeToDomainName: strings.Repeat(" ", len(schema.Domain)),
		ImportsModels:               importsModels,
		ImportsRepository:           importsRepository,
		ImportsUsecase:              importsUsecase,
		ImportsGrpcDelivery:         importsGrpcDelivery,
		Enums:                       enums,
		Entities:                    entities,
		TypesRepository:             typesRepository,
		TypesUsecase:                typesUsecase,
		MethodsRepository:           typeParser.MethodsRepository,
		MethodsUsecase:              typeParser.MethodsUsecase,
		MethodsGrpcDelivery:         typeParser.MethodsGrpcDelivery,
	}

	templateManager := template.NewTemplateManager()
	for k, v := range templatesNamesValues {
		err := templateManager.AddTemplate(k, v)
		if err != nil {
			return nil, err
		}
	}

	files := []*File{}
	domainKebab := formatter.PascalToKebab(schema.Domain)

	if len(templInput.Enums) > 0 || len(templInput.Entities) > 0 {
		models, err := templateManager.Parse("models", templInput)
		if err != nil {
			return nil, err
		}
		files = append(files, &File{
			Name:      "models/" + domainKebab + ".go",
			Content:   models,
			Overwrite: true,
		})
	}
	if len(templInput.MethodsRepository) > 0 {
		repository, err := templateManager.Parse("repository", templInput)
		if err != nil {
			return nil, err
		}
		files = append(files, &File{
			Name:      fmt.Sprintf("repository/%s/%s.go", domainKebab, domainKebab),
			Content:   repository,
			Overwrite: true,
		})

		repositoryStruct, err := templateManager.Parse("repository-struct", templInput)
		if err != nil {
			return nil, err
		}
		files = append(files, &File{
			Name:    fmt.Sprintf("repository/%s/implementation.go", domainKebab),
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
				Name:    fmt.Sprintf("repository/%s/%s.go", domainKebab, methodNameKebab),
				Content: repositoryMethod,
			})
		}
	}
	if len(templInput.MethodsUsecase) > 0 {
		usecase, err := templateManager.Parse("usecase", templInput)
		if err != nil {
			return nil, err
		}
		files = append(files, &File{
			Name:      fmt.Sprintf("usecase/%s/%s.go", domainKebab, domainKebab),
			Content:   usecase,
			Overwrite: true,
		})

		usecaseStruct, err := templateManager.Parse("usecase-struct", templInput)
		if err != nil {
			return nil, err
		}
		files = append(files, &File{
			Name:    fmt.Sprintf("usecase/%s/implementation.go", domainKebab),
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
				Name:    fmt.Sprintf("usecase/%s/%s.go", domainKebab, methodNameKebab),
				Content: usecaseMethod,
			})
		}
	}
	if len(templInput.MethodsGrpcDelivery) > 0 {
		grpcModule, err := templateManager.Parse("grpc-delivery-module", templInput)
		if err != nil {
			return nil, err
		}
		files = append(
			files,
			&File{
				Name:      fmt.Sprintf("delivery/grpc/%s/%s.go", templInput.DomainSnake, domainKebab),
				Content:   grpcModule,
				Overwrite: true,
			},
			&File{
				Name:    "delivery/grpc/grpc.go",
				Content: templates.GrpcDelivery,
			},
		)
	}

	files = append(
		files,
		&File{
			Name:      "adapters/validator.go",
			Content:   templates.ValidatorTempl,
			Overwrite: true,
		},
		&File{
			Name:    "adapters/go-validator/go-validator.go",
			Content: templates.ValidatorImplementationTempl,
		},
	)

	files = append(
		files,
		&File{
			Name:      "delivery/delivery.go",
			Content:   templates.DeliveryTempl,
			Overwrite: true,
		},
		&File{
			Name:      "delivery/gratefully-shutdown.go",
			Content:   templates.GratefullyShutdownTempl,
			Overwrite: true,
		},
	)

	files = append(
		files,
		&File{
			Name:      "utils/sync.go",
			Content:   templates.UtilsSyncTempl,
			Overwrite: true,
		},
	)

	files = append(
		files,
		&File{
			Name:    "cmd/main.go",
			Content: templates.MainTempl,
		},
	)

	return files, nil
}
