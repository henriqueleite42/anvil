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

func Parse(schema *schemas.Schema, silent bool, outputFolderPath string) error {
	if schema.Domain == "" {
		return fmt.Errorf("no domain specified")
	}
	if schema.Delivery == nil {
		return fmt.Errorf("no delivery specified")
	}
	if schema.Delivery.Grpc == nil {
		return fmt.Errorf("no gRPC delivery specified")
	}
	if schema.Delivery.Grpc.Rpcs == nil {
		return fmt.Errorf("no RPCs specified for gRPC delivery")
	}
	if schema.Usecase == nil {
		return fmt.Errorf("no usecases to deliver")
	}
	if schema.Usecase.Methods == nil || schema.Usecase.Methods.Methods == nil {
		return fmt.Errorf("no usecases methods to deliver")
	}

	rpcs := make([]*schemas.DeliveryGrpcRpc, 0, len(schema.Delivery.Grpc.Rpcs))
	for _, v := range schema.Delivery.Grpc.Rpcs {
		rpcs = append(rpcs, v)
	}
	sort.Slice(rpcs, func(i, j int) bool {
		return rpcs[i].Order < rpcs[j].Order
	})

	templateManager := template.NewTemplateManager()
	for k, v := range templatesNamesValues {
		err := templateManager.AddTemplate(k, v)
		if err != nil {
			return err
		}
	}

	contractGoTypesParser, err := types_parser.NewTypeParser(schema)
	if err != nil {
		return err
	}
	implementationGoTypesParser, err := types_parser.NewTypeParser(schema)
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
		method, ok := schema.Usecase.Methods.Methods[v.UsecaseMethodHash]
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

			_, err := contractGoTypesParser.ParseType(inputType, nil)
			if err != nil {
				return err
			}

			t, err := grpcParser.GoToProto(&grpc.GoToProtoInput{
				Type:                    inputType,
				MethodName:              methodName,
				VariableName:            "i",
				PrefixForVariableNaming: "Input",
				HasOutput:               true,
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

			_, err := contractGoTypesParser.ParseType(outputType, nil)
			if err != nil {
				return err
			}

			t, err := grpcParser.ProtoToGo(&grpc.ProtoToGoInput{
				Type:                    outputType,
				MethodName:              methodName,
				VariableName:            "result",
				PrefixForVariableNaming: "Output",
				HasOutput:               true,
			})
			if err != nil {
				return err
			}

			output = t
		} else {
			implementationGoTypesParser.AddImport("google.golang.org/protobuf/types/known/emptypb")
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
	types := contractGoTypesParser.GetMapTypes()

	// -----------------------------
	//
	// build template
	//
	// -----------------------------

	templInput := &templates.TemplInput{
		Domain:                      schema.Domain,
		DomainCamel:                 formatter.PascalToCamel(schema.Domain),
		DomainSnake:                 formatter.PascalToSnake(schema.Domain),
		SpacingRelativeToDomainName: strings.Repeat(" ", len(schema.Domain)),
		ImportsContract:             importsContract,
		ImportsImplementation:       importsImplementation,
		Enums:                       enums,
		Types:                       types,
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

	err = WriteFile(schema.Domain, outputFolderPath, "contract", contract)
	if err != nil {
		return err
	}

	err = WriteFile(schema.Domain, outputFolderPath, "implementation", implementation)
	if err != nil {
		return err
	}

	return nil
}
