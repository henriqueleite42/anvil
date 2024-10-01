package internal

import (
	"fmt"
	"sort"

	"github.com/henriqueleite42/anvil/generators/grpc/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	"github.com/henriqueleite42/anvil/language-helpers/golang/template"
)

type parser struct {
	schema                  *schemas.Schema
	imports                 map[string]bool
	methods                 []*templates.ProtofileTemplInputMethod
	enums                   map[string]*templates.ProtofileTemplInputEnum
	typesToAvoidDuplication map[string]*templates.ProtofileTemplInputType
	types                   []*templates.ProtofileTemplInputType
}

var templatesNamesValues = map[string]string{
	"protofile": templates.ProtofileTempl,
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

	// -----------------------------
	//
	// Parse methods
	//
	// -----------------------------

	parserInstance := &parser{
		schema:                  schema,
		imports:                 map[string]bool{},
		methods:                 make([]*templates.ProtofileTemplInputMethod, 0, len(rpcs)),
		enums:                   map[string]*templates.ProtofileTemplInputEnum{},
		typesToAvoidDuplication: map[string]*templates.ProtofileTemplInputType{},
		types:                   []*templates.ProtofileTemplInputType{},
	}

	emptyMsg := "google.protobuf.Empty"

	for _, v := range rpcs {
		method, ok := schema.Usecase.Methods.Methods[v.UsecaseMethodHash]
		if !ok {
			return fmt.Errorf("usecase method \"%s\" not found", v.UsecaseMethodHash)
		}

		var input *string = nil
		if method.Input != nil {
			if method.Input.TypeHash == "" {
				return fmt.Errorf("missing \"TypeHash\" for usecase method \"%s\"", method.Name)
			}

			inputType, ok := schema.Types.Types[method.Input.TypeHash]
			if !ok {
				return fmt.Errorf("type \"%s\" not found for usecase method \"%s\"", method.Input.TypeHash, method.Name)
			}

			inputTypeResolved, err := parserInstance.resolveType(inputType)
			if err != nil {
				return err
			}

			input = &inputTypeResolved.Name
		} else {
			parserInstance.imports["google/protobuf/empty.proto"] = true
			input = &emptyMsg
		}

		var output *string = nil
		if method.Output != nil {
			if method.Output.TypeHash == "" {
				return fmt.Errorf("missing \"TypeHash\" for usecase method \"%s\"", method.Name)
			}

			outputType, ok := schema.Types.Types[method.Output.TypeHash]
			if !ok {
				return fmt.Errorf("type \"%s\" not found for usecase method \"%s\"", method.Output.TypeHash, method.Name)
			}

			outputTypeResolved, err := parserInstance.resolveType(outputType)
			if err != nil {
				return err
			}

			output = &outputTypeResolved.Name
		} else {
			parserInstance.imports["google/protobuf/empty.proto"] = true
			output = &emptyMsg
		}

		parserInstance.methods = append(parserInstance.methods, &templates.ProtofileTemplInputMethod{
			Name:   method.Name,
			Input:  input,
			Output: output,
		})
	}

	// -----------------------------
	//
	// Prepare template
	//
	// -----------------------------

	templInput := &templates.ProtofileTemplInput{
		Domain:  schema.Domain,
		Imports: make([]string, 0, len(parserInstance.imports)),
		Enums:   make([]*templates.ProtofileTemplInputEnum, 0, len(parserInstance.enums)),
		Methods: parserInstance.methods,
		Types:   parserInstance.types,
	}
	for k := range parserInstance.imports {
		templInput.Imports = append(templInput.Imports, k)
	}
	sort.Slice(templInput.Imports, func(i, j int) bool {
		return templInput.Imports[i] < templInput.Imports[j]
	})
	for _, v := range parserInstance.enums {
		templInput.Enums = append(templInput.Enums, v)
	}
	sort.Slice(templInput.Enums, func(i, j int) bool {
		return templInput.Enums[i].Name < templInput.Enums[j].Name
	})

	protofile, err := templateManager.Parse("protofile", templInput)
	if err != nil {
		return err
	}

	err = WriteFile(schema, outputFolderPath, protofile)
	if err != nil {
		return err
	}

	return nil
}
