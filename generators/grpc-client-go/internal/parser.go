package internal

import (
	"fmt"
	"sort"
	"strings"

	"github.com/henriqueleite42/anvil/generators/grpc-client-go/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/formatter"
	"github.com/henriqueleite42/anvil/language-helpers/golang/grpc"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	"github.com/henriqueleite42/anvil/language-helpers/golang/template"
	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

var templatesNamesValues = map[string]string{
	"contract":       templates.ContractTempl,
	"implementation": templates.ImplementationTempl,
}

func Parse(schema *schemas.AnvpSchema, silent bool, outputFolderPath *string) error {
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

	for curDomain, v := range schema.Deliveries.Deliveries {
		if v.Grpc == nil || v.Grpc.Rpcs == nil {
			continue
		}
		if _, ok := schema.Usecases.Usecases[curDomain]; !ok {
			return fmt.Errorf("no usecases specified for domain \"%s\"", curDomain)
		}

		domainSnake := formatter.PascalToSnake(curDomain)

		rpcs := make([]*schemas.DeliveryGrpcRpc, 0, len(v.Grpc.Rpcs))
		for _, v := range v.Grpc.Rpcs {
			rpcs = append(rpcs, v)
		}
		sort.Slice(rpcs, func(i, j int) bool {
			return rpcs[i].Order < rpcs[j].Order
		})

		pkg := domainSnake + "_grpc_client"

		contractGoTypesParser, err := types_parser.NewTypeParser(&types_parser.NewTypeParserInput{
			Schema:        schema,
			EnumsPkg:      pkg,
			TypesPkg:      pkg,
			EventsPkg:     pkg,
			EntitiesPkg:   pkg,
			RepositoryPkg: pkg,
			UsecasePkg:    pkg,
		})
		if err != nil {
			return err
		}
		implementationGoTypesParser, err := types_parser.NewTypeParser(&types_parser.NewTypeParserInput{
			Schema:        schema,
			EnumsPkg:      pkg,
			TypesPkg:      pkg,
			EventsPkg:     pkg,
			EntitiesPkg:   pkg,
			RepositoryPkg: pkg,
			UsecasePkg:    pkg,
		})
		if err != nil {
			return err
		}

		implementationGoTypesParser.AddImport("time")
		implementationGoTypesParser.AddImport("context")
		implementationGoTypesParser.AddImport("errors")
		implementationGoTypesParser.AddImport("google.golang.org/grpc")
		implementationGoTypesParser.AddImport("google.golang.org/grpc/credentials/insecure")

		grpcParser := grpc.NewGrpcParser(&grpc.NewGrpcParserInput{
			Schema:       schema,
			GoTypeParser: implementationGoTypesParser,
		})

		// -----------------------------
		//
		// Parse methods
		//
		// -----------------------------

		methods := []*templates.TemplMethod{}
		for _, v := range rpcs {
			method, ok := schema.Usecases.Usecases[curDomain].Methods.Methods[v.UsecaseMethodHash]
			if !ok {
				return fmt.Errorf("usecase method \"%s\" not found", v.UsecaseMethodHash)
			}

			methodName := method.Name

			var input *grpc.Type = nil
			if method.Input != nil {
				if method.Input.TypeHash == "" {
					return fmt.Errorf("missing \"TypeHash\" for usecase method \"%s\"", method.Name)
				}

				inputType, ok := schema.Types.Types[method.Input.TypeHash]
				if !ok {
					return fmt.Errorf("type \"%s\" not found for usecase method \"%s\"", method.Input.TypeHash, method.Name)
				}

				_, err := contractGoTypesParser.ParseType(inputType)
				if err != nil {
					return err
				}

				t, err := grpcParser.GoToProto(&grpc.GoToProtoInput{
					Type:                     inputType,
					MethodName:               methodName,
					VariableToAccessTheValue: "i",
					HasOutput:                method.Output != nil,
					CurPkg:                   pkg,
				})
				if err != nil {
					return err
				}

				input = t
			} else {
				implementationGoTypesParser.AddImport("google.golang.org/protobuf/types/known/emptypb")
			}

			var output *grpc.Type = nil
			if method.Output != nil {
				if method.Output.TypeHash == "" {
					return fmt.Errorf("missing \"TypeHash\" for usecase method \"%s\"", method.Name)
				}

				outputType, ok := schema.Types.Types[method.Output.TypeHash]
				if !ok {
					return fmt.Errorf("type \"%s\" not found for usecase method \"%s\"", method.Output.TypeHash, method.Name)
				}

				_, err := contractGoTypesParser.ParseType(outputType)
				if err != nil {
					return err
				}

				t, err := grpcParser.ProtoToGo(&grpc.ProtoToGoInput{
					Type:                     outputType,
					MethodName:               methodName,
					VariableToAccessTheValue: "result",
					HasOutput:                true,
					CurPkg:                   pkg,
				})
				if err != nil {
					return err
				}

				output = t
			}

			methods = append(methods, &templates.TemplMethod{
				MethodName:      methodName,
				MethodNameCamel: formatter.PascalToCamel(method.Name),
				Input:           input,
				Output:          output,
			})
		}

		// -----------------------------
		//
		// prepare values
		//
		// -----------------------------

		contractGoTypesParser.AddImport("time")
		importsContract := contractGoTypesParser.GetImports()

		importsImplementation := implementationGoTypesParser.GetImports()

		enums := contractGoTypesParser.GetEnums()
		types := contractGoTypesParser.GetTypes()
		entities := contractGoTypesParser.GetEntities()
		events := contractGoTypesParser.GetEvents()
		repository := contractGoTypesParser.GetRepository()
		usecase := contractGoTypesParser.GetUsecase()

		allTypes := make(
			[]*types_parser.Type,
			0,
			len(types)+len(entities)+len(events)+len(repository)+len(usecase),
		)
		allTypes = append(allTypes, types...)
		allTypes = append(allTypes, entities...)
		allTypes = append(allTypes, events...)
		allTypes = append(allTypes, repository...)
		allTypes = append(allTypes, usecase...)

		// -----------------------------
		//
		// build template
		//
		// -----------------------------

		templInput := &templates.TemplInput{
			Domain:                      curDomain,
			DomainCamel:                 formatter.PascalToCamel(curDomain),
			DomainSnake:                 domainSnake,
			SpacingRelativeToDomainName: strings.Repeat(" ", len(curDomain)),
			ImportsContract:             importsContract,
			ImportsImplementation:       importsImplementation,
			Enums:                       enums,
			Types:                       allTypes,
			Methods:                     methods,
		}

		contract, err := templateManager.Parse("contract", templInput)
		if err != nil {
			return err
		}
		implementation, err := templateManager.Parse("implementation", templInput)
		if err != nil {
			return err
		}

		err = WriteFile(curDomain, outputFolderPath, "contract", contract)
		if err != nil {
			return err
		}

		err = WriteFile(curDomain, outputFolderPath, "implementation", implementation)
		if err != nil {
			return err
		}
	}

	return nil
}
