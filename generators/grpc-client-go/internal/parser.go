package internal

import (
	"fmt"
	"sort"
	"strings"

	generator_config "github.com/henriqueleite42/anvil/generators/grpc-client-go/config"
	"github.com/henriqueleite42/anvil/generators/grpc-client-go/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/formatter"
	"github.com/henriqueleite42/anvil/language-helpers/golang/grpc"
	"github.com/henriqueleite42/anvil/language-helpers/golang/imports"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	"github.com/henriqueleite42/anvil/language-helpers/golang/template"
	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

var templatesNamesValues = map[string]string{
	"contract":       templates.ContractTempl,
	"implementation": templates.ImplementationTempl,
}

func Parse(schema *schemas.AnvpSchema, config *generator_config.GeneratorConfig, silent bool) error {
	if schema.Deliveries == nil || schema.Deliveries.Deliveries == nil {
		return fmt.Errorf("no deliveries specified")
	}
	if schema.Usecases == nil || schema.Usecases.Usecases == nil {
		return fmt.Errorf("no usecases specified")
	}

	templateManager := template.NewTemplateManager()
	for k, v := range templatesNamesValues {
		err := templateManager.AddTemplate(k, v)
		if err != nil {
			return err
		}
	}

	resolveEnumImpt := func(e *schemas.Enum) *imports.Import {
		domainSnake := formatter.PascalToSnake(e.Domain)
		pkg := domainSnake + "_grpc_client"
		return imports.NewImport(config.ClientsModuleName+"/"+domainSnake, &pkg)
	}
	resolveTypeImpt := func(t *schemas.Type) *imports.Import {
		domainSnake := formatter.PascalToSnake(t.Domain)
		pkg := domainSnake + "_grpc_client"
		return imports.NewImport(config.ClientsModuleName+"/"+domainSnake, &pkg)
	}

	pbImport := imports.NewImport(config.ClientsModuleName+"/pb", nil)

	goTypesParser, err := types_parser.NewTypeParser(&types_parser.NewTypeParserInput{
		Schema: schema,

		GetEnumsImport:      resolveEnumImpt,
		GetTypesImport:      resolveTypeImpt,
		GetEventsImport:     resolveTypeImpt,
		GetEntitiesImport:   resolveTypeImpt,
		GetRepositoryImport: resolveTypeImpt,
		GetUsecaseImport:    resolveTypeImpt,
	})
	if err != nil {
		return err
	}

	grpcParser := grpc.NewGrpcParser(&grpc.NewGrpcParserInput{
		Schema:                schema,
		GoTypeParser:          goTypesParser,
		GetEnumConversionImpt: resolveEnumImpt,
	})

	contractsImportsPerDomain := make(map[string]imports.ImportsManager, len(schema.Schemas))
	implementationImportsPerDomain := make(map[string]imports.ImportsManager, len(schema.Schemas))
	methodsPerDomain := make(map[string][]*templates.TemplMethod, len(schema.Schemas))

	for _, shm := range schema.Schemas {
		contractsImportsPerDomain[shm.Domain] = imports.NewImportsManager()

		implementationImportsPerDomain[shm.Domain] = imports.NewImportsManager()
		implementationImportsPerDomain[shm.Domain].MergeImport(pbImport)
		implementationImportsPerDomain[shm.Domain].AddImport("time", nil)
		implementationImportsPerDomain[shm.Domain].AddImport("context", nil)
		implementationImportsPerDomain[shm.Domain].AddImport("errors", nil)
		implementationImportsPerDomain[shm.Domain].AddImport("google.golang.org/grpc", nil)
		implementationImportsPerDomain[shm.Domain].AddImport("google.golang.org/grpc/credentials/insecure", nil)
	}

	if schema.Enums != nil && schema.Enums.Enums != nil {
		for _, e := range schema.Enums.Enums {
			_, err := goTypesParser.ParseEnum(e)
			if err != nil {
				return err
			}
		}
	}

	if schema.Deliveries != nil && schema.Deliveries.Deliveries != nil {
		for curDomain, v := range schema.Deliveries.Deliveries {
			if v.Grpc == nil || v.Grpc.Rpcs == nil {
				continue
			}
			if _, ok := schema.Usecases.Usecases[curDomain]; !ok {
				return fmt.Errorf("no usecases specified for domain \"%s\"", curDomain)
			}

			rpcs := make([]*schemas.DeliveryGrpcRpc, 0, len(v.Grpc.Rpcs))
			for _, v := range v.Grpc.Rpcs {
				rpcs = append(rpcs, v)
			}
			sort.Slice(rpcs, func(i, j int) bool {
				return rpcs[i].Order < rpcs[j].Order
			})

			// -----------------------------
			//
			// Parse methods
			//
			// -----------------------------

			methodsPerDomain[curDomain] = make([]*templates.TemplMethod, 0, len(rpcs))
			for _, v := range rpcs {
				method, ok := schema.Usecases.Usecases[curDomain].Methods.Methods[v.UsecaseMethodHash]
				if !ok {
					return fmt.Errorf("usecase method \"%s\" not found", v.UsecaseMethodHash)
				}

				methodName := method.Name

				var input *grpc.ConvertedValue = nil
				if method.Input != nil {
					if method.Input.TypeHash == "" {
						return fmt.Errorf("missing \"TypeHash\" for usecase method \"%s\"", method.Name)
					}

					inputType, ok := schema.Types.Types[method.Input.TypeHash]
					if !ok {
						return fmt.Errorf("type \"%s\" not found for usecase method \"%s\"", method.Input.TypeHash, method.Name)
					}

					_, err := goTypesParser.ParseType(inputType)
					if err != nil {
						return err
					}

					t, err := grpcParser.GoToProto(&grpc.ConverterInput{
						Type:            inputType,
						CurModuleImport: resolveTypeImpt(inputType),
						PbModuleImport:  pbImport,
						VarToConvert:    "i",
					})
					if err != nil {
						return err
					}

					implementationImportsPerDomain[curDomain].MergeImports(t.ImportsUnorganized)

					input = t
				} else {
					implementationImportsPerDomain[curDomain].AddImport("google.golang.org/protobuf/types/known/emptypb", nil)
				}

				var output *grpc.ConvertedValue = nil
				if method.Output != nil {
					if method.Output.TypeHash == "" {
						return fmt.Errorf("missing \"TypeHash\" for usecase method \"%s\"", method.Name)
					}

					outputType, ok := schema.Types.Types[method.Output.TypeHash]
					if !ok {
						return fmt.Errorf("type \"%s\" not found for usecase method \"%s\"", method.Output.TypeHash, method.Name)
					}

					_, err := goTypesParser.ParseType(outputType)
					if err != nil {
						return err
					}

					t, err := grpcParser.ProtoToGo(&grpc.ConverterInput{
						Type:            outputType,
						CurModuleImport: resolveTypeImpt(outputType),
						PbModuleImport:  pbImport,
						VarToConvert:    "result",
					})
					if err != nil {
						return err
					}

					implementationImportsPerDomain[curDomain].MergeImports(t.ImportsUnorganized)

					output = t
				}

				methodsPerDomain[curDomain] = append(methodsPerDomain[curDomain], &templates.TemplMethod{
					MethodName:      methodName,
					MethodNameCamel: formatter.PascalToCamel(method.Name),
					Input:           input,
					Output:          output,
				})
			}
		}
	}

	enums := goTypesParser.GetEnums()
	types := goTypesParser.GetTypes()
	entities := goTypesParser.GetEntities()
	events := goTypesParser.GetEvents()
	repository := goTypesParser.GetRepository()
	usecase := goTypesParser.GetUsecase()

	enumsPerDomain := map[string][]*templates.TemplEnum{}
	typesPerDomain := map[string][]*templates.TemplType{}

	for _, e := range enums {
		if _, ok := enumsPerDomain[e.AnvilEnum.Domain]; !ok {
			enumsPerDomain[e.AnvilEnum.Domain] = []*templates.TemplEnum{}
		}
		enumsPerDomain[e.AnvilEnum.Domain] = append(enumsPerDomain[e.AnvilEnum.Domain], enumToTemplEnum(e))
	}
	for _, t := range types {
		if _, ok := typesPerDomain[t.AnvilType.Domain]; !ok {
			typesPerDomain[t.AnvilType.Domain] = []*templates.TemplType{}
		}
		contractsImportsPerDomain[t.AnvilType.Domain].MergeImports(t.GetImportsUnorganized())
		pkg := resolveTypeImpt(t.AnvilType)
		typesPerDomain[t.AnvilType.Domain] = append(typesPerDomain[t.AnvilType.Domain], typeToTemplType(pkg.Alias, t))
	}
	for _, t := range entities {
		if _, ok := typesPerDomain[t.AnvilType.Domain]; !ok {
			typesPerDomain[t.AnvilType.Domain] = []*templates.TemplType{}
		}
		contractsImportsPerDomain[t.AnvilType.Domain].MergeImports(t.GetImportsUnorganized())
		pkg := resolveTypeImpt(t.AnvilType)
		typesPerDomain[t.AnvilType.Domain] = append(typesPerDomain[t.AnvilType.Domain], typeToTemplType(pkg.Alias, t))
	}
	for _, t := range events {
		if _, ok := typesPerDomain[t.AnvilType.Domain]; !ok {
			typesPerDomain[t.AnvilType.Domain] = []*templates.TemplType{}
		}
		contractsImportsPerDomain[t.AnvilType.Domain].MergeImports(t.GetImportsUnorganized())
		pkg := resolveTypeImpt(t.AnvilType)
		typesPerDomain[t.AnvilType.Domain] = append(typesPerDomain[t.AnvilType.Domain], typeToTemplType(pkg.Alias, t))
	}
	for _, t := range repository {
		if _, ok := typesPerDomain[t.AnvilType.Domain]; !ok {
			typesPerDomain[t.AnvilType.Domain] = []*templates.TemplType{}
		}
		contractsImportsPerDomain[t.AnvilType.Domain].MergeImports(t.GetImportsUnorganized())
		pkg := resolveTypeImpt(t.AnvilType)
		typesPerDomain[t.AnvilType.Domain] = append(typesPerDomain[t.AnvilType.Domain], typeToTemplType(pkg.Alias, t))
	}
	for _, t := range usecase {
		if _, ok := typesPerDomain[t.AnvilType.Domain]; !ok {
			typesPerDomain[t.AnvilType.Domain] = []*templates.TemplType{}
		}
		contractsImportsPerDomain[t.AnvilType.Domain].MergeImports(t.GetImportsUnorganized())
		pkg := resolveTypeImpt(t.AnvilType)
		typesPerDomain[t.AnvilType.Domain] = append(typesPerDomain[t.AnvilType.Domain], typeToTemplType(pkg.Alias, t))
	}

	// -----------------------------
	//
	// build template
	//
	// -----------------------------
	for _, shm := range schema.Schemas {
		domainSnake := formatter.PascalToSnake(shm.Domain)
		pkg := domainSnake + "_grpc_client"

		methods, hasMethods := methodsPerDomain[shm.Domain]

		var importsImplementation [][]string = nil
		if im, ok := implementationImportsPerDomain[shm.Domain]; ok {
			importsImplementation = im.ResolveImports(pkg)
		}
		if hasMethods {
			contractsImportsPerDomain[shm.Domain].AddImport("time", nil)
		}
		importsContract := contractsImportsPerDomain[shm.Domain].ResolveImports(pkg)

		templInput := &templates.TemplInput{
			Domain:                      shm.Domain,
			DomainCamel:                 formatter.PascalToCamel(shm.Domain),
			DomainSnake:                 domainSnake,
			SpacingRelativeToDomainName: strings.Repeat(" ", len(shm.Domain)),
			ImportsContract:             importsContract,
			ImportsImplementation:       importsImplementation,
			Enums:                       enumsPerDomain[shm.Domain],
			Types:                       typesPerDomain[shm.Domain],
			Methods:                     methods,
		}

		contract, err := templateManager.Parse("contract", templInput)
		if err != nil {
			return err
		}
		err = WriteFile(shm.Domain, config.OutDir, "contract", contract)
		if err != nil {
			return err
		}

		if hasMethods {
			implementation, err := templateManager.Parse("implementation", templInput)
			if err != nil {
				return err
			}
			err = WriteFile(shm.Domain, config.OutDir, "implementation", implementation)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
