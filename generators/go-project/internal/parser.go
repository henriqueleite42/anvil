package internal

import (
	"fmt"
	"sort"
	"strings"

	"github.com/henriqueleite42/anvil/cli/formatter"
	"github.com/henriqueleite42/anvil/cli/schemas"
	"github.com/henriqueleite42/anvil/cli/template"
	"github.com/henriqueleite42/anvil/generators/go-project/internal/parser"
	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
)

var templatesNamesValues = map[string]string{
	"models":            templates.ModelsTempl,
	"repository":        templates.RepositoryTempl,
	"repository-struct": templates.RepositoryStructTempl,
	"repository-method": templates.RepositoryMethodTempl,
	"usecase":           templates.UsecaseTempl,
	"usecase-struct":    templates.UsecaseStructTempl,
	"usecase-method":    templates.UsecaseMethodTempl,
}

type File struct {
	Name    string
	Content string
}

func getImports(imports map[string]bool) []string {
	importsStd := make([]string, 0, len(imports))
	importsExt := make([]string, 0, len(imports))
	for k := range imports {
		impt := fmt.Sprintf("	\"%s\"", k)
		parts := strings.Split(k, "/")
		if strings.Contains(parts[0], ".") {
			importsExt = append(importsExt, impt)
		} else {
			importsStd = append(importsStd, impt)
		}
	}
	sort.Slice(importsStd, func(i, j int) bool {
		return importsStd[i] < importsStd[j]
	})
	sort.Slice(importsExt, func(i, j int) bool {
		return importsExt[i] < importsExt[j]
	})
	importsResolved := make([]string, 0, len(imports)+1)
	if len(importsStd) > 0 {
		importsResolved = append(importsResolved, importsStd...)
	}
	if len(importsStd) > 0 && len(importsExt) > 0 {
		importsResolved = append(importsResolved, "")
	}
	if len(importsExt) > 0 {
		importsResolved = append(importsResolved, importsExt...)
	}
	return importsResolved
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

	typeParser := &parser.Parser{
		ModelsPath:    "foo/models",
		ModelsPkgName: "models",
		Schema:        schema,

		ImportsModels:                     map[string]bool{},
		ImportsRepository:                 map[string]bool{},
		ImportsUsecase:                    map[string]bool{},
		Enums:                             make(map[string]*templates.TemplEnum, lenEnums),
		Entities:                          make([]*templates.TemplType, 0, lenEntities),
		TypesRepository:                   []*templates.TemplType{},
		TypesRepositoryToAvoidDuplication: map[string]*templates.TemplType{}, // TODO improve this, find a way to merge both
		TypesUsecase:                      []*templates.TemplType{},
		TypesUsecaseToAvoidDuplication:    map[string]*templates.TemplType{},
		MethodsRepository:                 make([]*templates.TemplMethod, 0, lenRepository),
		MethodsUsecase:                    make([]*templates.TemplMethod, 0, lenUsecase),
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
			var inputTypeHash string
			if v.Input != nil {
				inputTypeHash = v.Input.TypeHash
			}
			var outputTypeHash string
			if v.Output != nil {
				outputTypeHash = v.Output.TypeHash
			}

			err := typeParser.ResolveMethod(parser.Kind_Repository, v.Name, inputTypeHash, outputTypeHash)
			if err != nil {
				return nil, err
			}
		}
	}

	if lenUsecase > 0 {
		for _, v := range schema.Usecase.Methods.Methods {
			var inputTypeHash string
			if v.Input != nil {
				inputTypeHash = v.Input.TypeHash
			}
			var outputTypeHash string
			if v.Output != nil {
				outputTypeHash = v.Output.TypeHash
			}

			err := typeParser.ResolveMethod(parser.Kind_Usecase, v.Name, inputTypeHash, outputTypeHash)
			if err != nil {
				return nil, err
			}
		}
	}

	importsModels := getImports(typeParser.ImportsModels)
	importsRepository := getImports(typeParser.ImportsRepository)
	importsUsecase := getImports(typeParser.ImportsUsecase)

	enums := make([]*templates.TemplEnum, 0, len(typeParser.Enums))
	for _, v := range typeParser.Enums {
		enums = append(enums, v)
	}
	sort.Slice(enums, func(i, j int) bool {
		return enums[i].Name < enums[j].Name
	})

	entities := make([]*templates.TemplType, 0, len(typeParser.Entities))
	for _, v := range typeParser.Entities {
		entities = append(entities, v)
	}
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Name < entities[j].Name
	})

	templInput := &templates.TemplInput{
		Domain:            schema.Domain,
		DomainSnake:       formatter.PascalToSnake(schema.Domain),
		ImportsModels:     importsModels,
		ImportsRepository: importsRepository,
		ImportsUsecase:    importsUsecase,
		Enums:             enums,
		Entities:          entities,
		TypesRepository:   typeParser.TypesRepository,
		TypesUsecase:      typeParser.TypesUsecase,
		MethodsRepository: typeParser.MethodsRepository,
		MethodsUsecase:    typeParser.MethodsUsecase,
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
			Name:    "models/" + domainKebab + ".go",
			Content: models,
		})
	}
	if len(templInput.MethodsRepository) > 0 {
		repository, err := templateManager.Parse("repository", templInput)
		if err != nil {
			return nil, err
		}
		files = append(files, &File{
			Name:    fmt.Sprintf("repository/%s/%s.go", domainKebab, domainKebab),
			Content: repository,
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
			Name:    fmt.Sprintf("usecase/%s/%s.go", domainKebab, domainKebab),
			Content: usecase,
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

	return files, nil
}
