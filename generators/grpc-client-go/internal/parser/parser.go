package parser

import (
	"fmt"
	"sort"
	"strings"

	"github.com/henriqueleite42/anvil/cli/formatter"
	"github.com/henriqueleite42/anvil/cli/schemas"
	"github.com/henriqueleite42/anvil/cli/template"
	"github.com/henriqueleite42/anvil/generators/grpc-client-go/internal/templates"
)

type parserManager struct {
	schema                *schemas.Schema
	templateManager       template.TemplateManager
	importsContract       map[string]bool
	importsImplementation map[string]bool
	enums                 map[string]*templates.TemplEnum
	types                 map[string]*templates.TemplType
	methods               []*templates.TemplMethod
}

var templatesNamesValues = map[string]string{
	"contract":            templates.ContractTempl,
	"implementation":      templates.ImplementationTempl,
	"input-prop-list":     templates.InputPropListTempl,
	"input-prop-map":      templates.InputPropMapTempl,
	"input-prop-optional": templates.InputPropOptionalTempl,
}

func Parse(schema *schemas.Schema) (string, string, error) {
	if schema.Domain == "" {
		return "", "", fmt.Errorf("no domain specified")
	}
	if schema.Delivery == nil {
		return "", "", fmt.Errorf("no delivery specified")
	}
	if schema.Delivery.Grpc == nil {
		return "", "", fmt.Errorf("no gRPC delivery specified")
	}
	if schema.Delivery.Grpc.Rpcs == nil {
		return "", "", fmt.Errorf("no RPCs specified for gRPC delivery")
	}
	if schema.Usecase == nil {
		return "", "", fmt.Errorf("no usecases to deliver")
	}
	if schema.Usecase.Methods == nil || schema.Usecase.Methods.Methods == nil {
		return "", "", fmt.Errorf("no usecases methods to deliver")
	}

	rpcs := []*schemas.DeliveryGrpcRpc{}
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
			return "", "", err
		}
	}

	parser := &parserManager{
		schema:          schema,
		templateManager: templateManager,
		importsContract: map[string]bool{
			"time": true,
		},
		importsImplementation: map[string]bool{
			"time":                   true,
			"context":                true,
			"errors":                 true,
			"google.golang.org/grpc": true,
			"google.golang.org/grpc/credentials/insecure": true,
		},
		enums:   map[string]*templates.TemplEnum{},
		types:   map[string]*templates.TemplType{},
		methods: []*templates.TemplMethod{},
	}

	// -----------------------------
	//
	// parse methods
	//
	// -----------------------------

	methods := []*templates.TemplMethod{}
	for _, v := range rpcs {
		method, ok := schema.Usecase.Methods.Methods[v.UsecaseMethodHash]
		if !ok {
			return "", "", fmt.Errorf("usecase method \"%s\" not found", v.UsecaseMethodHash)
		}

		methodName := method.Name

		var input []*templates.TemplMethodProp = nil
		var inputPropsPrepare []string = nil
		if method.Input != nil {
			if method.Input.TypeHash == "" {
				return "", "", fmt.Errorf("missing \"TypeHash\" for usecase method \"%s\"", method.Name)
			}

			inputType, ok := schema.Types.Types[method.Input.TypeHash]
			if !ok {
				return "", "", fmt.Errorf("type \"%s\" not found for usecase method \"%s\"", method.Input.TypeHash, method.Name)
			}

			props, propsPrepare, err := parser.toMethodInput(methodName, inputType)
			if err != nil {
				return "", "", err
			}

			input = props
			inputPropsPrepare = propsPrepare
		}

		var output []*templates.TemplMethodProp = nil
		var outputPropsPrepare []string = nil
		if method.Output != nil {
			if method.Output.TypeHash == "" {
				return "", "", fmt.Errorf("missing \"TypeHash\" for usecase method \"%s\"", method.Name)
			}

			outputType, ok := schema.Types.Types[method.Output.TypeHash]
			if !ok {
				return "", "", fmt.Errorf("type \"%s\" not found for usecase method \"%s\"", method.Output.TypeHash, method.Name)
			}

			props, propsPrepare, err := parser.toMethodOutput(outputType)
			if err != nil {
				return "", "", err
			}

			output = props
			outputPropsPrepare = propsPrepare
		}

		methods = append(methods, &templates.TemplMethod{
			MethodName:         methodName,
			MethodNameCamel:    formatter.PascalToCamel(method.Name),
			Input:              input,
			InputPropsPrepare:  inputPropsPrepare,
			Output:             output,
			OutputPropsPrepare: outputPropsPrepare,
		})
	}

	// -----------------------------
	//
	// prepare values
	//
	// -----------------------------

	importsContractStd := make([]string, 0, len(parser.importsContract))
	importsContractExt := make([]string, 0, len(parser.importsContract))
	for k := range parser.importsContract {
		impt := fmt.Sprintf("	\"%s\"", k)
		parts := strings.Split(k, "/")
		if strings.Contains(parts[0], ".") {
			importsContractExt = append(importsContractExt, impt)
		} else {
			importsContractStd = append(importsContractStd, impt)
		}
	}
	sort.Slice(importsContractStd, func(i, j int) bool {
		return importsContractStd[i] < importsContractStd[j]
	})
	sort.Slice(importsContractExt, func(i, j int) bool {
		return importsContractExt[i] < importsContractExt[j]
	})
	importsContract := make([]string, 0, len(parser.importsContract)+1)
	importsContract = append(importsContract, importsContractStd...)
	importsContract = append(importsContract, "")
	importsContract = append(importsContract, importsContractExt...)

	importsImplementationStd := make([]string, 0, len(parser.importsImplementation))
	importsImplementationExt := make([]string, 0, len(parser.importsImplementation))
	for k := range parser.importsImplementation {
		impt := fmt.Sprintf("	\"%s\"", k)
		parts := strings.Split(k, "/")
		if strings.Contains(parts[0], ".") {
			importsImplementationExt = append(importsImplementationExt, impt)
		} else {
			importsImplementationStd = append(importsImplementationStd, impt)
		}
	}
	sort.Slice(importsImplementationStd, func(i, j int) bool {
		return importsImplementationStd[i] < importsImplementationStd[j]
	})
	sort.Slice(importsImplementationExt, func(i, j int) bool {
		return importsImplementationExt[i] < importsImplementationExt[j]
	})
	importsImplementation := make([]string, 0, len(parser.importsImplementation)+1)
	importsImplementation = append(importsImplementation, importsImplementationStd...)
	importsImplementation = append(importsImplementation, "")
	importsImplementation = append(importsImplementation, importsImplementationExt...)

	enums := []*templates.TemplEnum{}
	for _, v := range parser.enums {
		enums = append(enums, v)
	}
	sort.Slice(enums, func(i, j int) bool {
		return enums[i].Name < enums[j].Name
	})
	types := []*templates.TemplType{}
	for _, v := range parser.types {
		types = append(types, v)
	}
	sort.Slice(types, func(i, j int) bool {
		return types[i].Name < types[j].Name
	})

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
		return "", "", err
	}
	implementation, err := templateManager.Parse("implementation", templInput)
	if err != nil {
		return "", "", err
	}

	return contract, implementation, nil
}
