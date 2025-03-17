package internal

import (
	"fmt"
	"sort"

	generator_config "github.com/henriqueleite42/anvil/generators/grpc/config"
	"github.com/henriqueleite42/anvil/generators/grpc/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/imports"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	"github.com/henriqueleite42/anvil/language-helpers/golang/template"
)

type parserManager struct {
	schema *schemas.AnvpSchema

	typesToAvoidDuplication map[string]*templates.ProtofileTemplInputType

	grpcTypesParser map[string]*grpcTypesParser
}

var templatesNamesValues = map[string]string{
	"protofile": templates.ProtofileTempl,
}

func Parse(schema *schemas.AnvpSchema, config *generator_config.GeneratorConfig, silent bool) error {
	if schema.Deliveries == nil || schema.Deliveries.Deliveries == nil {
		return fmt.Errorf("no delivery specified")
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

	parser := &parserManager{
		schema:                  schema,
		typesToAvoidDuplication: map[string]*templates.ProtofileTemplInputType{},
		grpcTypesParser:         make(map[string]*grpcTypesParser, len(schema.Schemas)),
	}

	for _, shm := range schema.Schemas {
		parser.grpcTypesParser[shm.Domain] = &grpcTypesParser{
			imports:  imports.NewImportsManager(),
			enums:    map[string]*templates.ProtofileTemplInputEnum{},
			types:    []*templates.ProtofileTemplInputType{},
			events:   []*templates.ProtofileTemplInputType{},
			entities: []*templates.ProtofileTemplInputType{},
		}
	}

	// -----------------------------
	//
	// Prepare types
	//
	// -----------------------------

	if schema.Enums != nil && schema.Enums.Enums != nil {
		for _, v := range schema.Enums.Enums {
			parser.resolveEnum(v)
		}
	}

	// -----------------------------
	//
	// Prepare rpcs
	//
	// -----------------------------

	for curDomain, v := range schema.Deliveries.Deliveries {
		emptyMsg := "google.protobuf.Empty"

		rpcs := make([]*schemas.DeliveryGrpcRpc, 0, len(v.Grpc.Rpcs))
		for _, v := range v.Grpc.Rpcs {
			rpcs = append(rpcs, v)
		}
		sort.Slice(rpcs, func(i, j int) bool {
			return rpcs[i].Order < rpcs[j].Order
		})

		parser.grpcTypesParser[curDomain].methods = make([]*templates.ProtofileTemplInputMethod, 0, len(rpcs))

		for _, v := range rpcs {
			method, ok := schema.Usecases.Usecases[curDomain].Methods.Methods[v.UsecaseMethodHash]
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

				inputTypeResolved, err := parser.resolveType(inputType, curDomain)
				if err != nil {
					return err
				}

				input = &inputTypeResolved.Name
			} else {
				parser.grpcTypesParser[curDomain].imports.AddImport("google/protobuf/empty.proto", nil)
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

				outputTypeResolved, err := parser.resolveType(outputType, curDomain)
				if err != nil {
					return err
				}

				output = &outputTypeResolved.Name
			} else {
				parser.grpcTypesParser[curDomain].imports.AddImport("google/protobuf/empty.proto", nil)
				output = &emptyMsg
			}

			parser.grpcTypesParser[curDomain].methods = append(parser.grpcTypesParser[curDomain].methods, &templates.ProtofileTemplInputMethod{
				Name:   method.Name,
				Input:  input,
				Output: output,
			})
		}
	}

	// -----------------------------
	//
	// Prepare templates
	//
	// -----------------------------

	for _, shm := range schema.Schemas {
		sort.Slice(parser.grpcTypesParser[shm.Domain].types, func(i, j int) bool {
			return parser.grpcTypesParser[shm.Domain].types[i].Name < parser.grpcTypesParser[shm.Domain].types[j].Name
		})
		templInput := &templates.ProtofileTemplInput{
			Domain:   shm.Domain,
			Imports:  make([]string, 0, parser.grpcTypesParser[shm.Domain].imports.GetImportsLen()),
			Enums:    make([]*templates.ProtofileTemplInputEnum, 0, len(parser.grpcTypesParser[shm.Domain].enums)),
			Methods:  parser.grpcTypesParser[shm.Domain].methods,
			Types:    parser.grpcTypesParser[shm.Domain].types,
			Events:   parser.grpcTypesParser[shm.Domain].events,
			Entities: parser.grpcTypesParser[shm.Domain].entities,
		}
		for _, v := range parser.grpcTypesParser[shm.Domain].imports.GetImportsUnorganized() {
			templInput.Imports = append(templInput.Imports, v.Path)
		}
		sort.Slice(templInput.Imports, func(i, j int) bool {
			return templInput.Imports[i] < templInput.Imports[j]
		})
		for _, v := range parser.grpcTypesParser[shm.Domain].enums {
			templInput.Enums = append(templInput.Enums, v)
		}
		sort.Slice(templInput.Enums, func(i, j int) bool {
			return templInput.Enums[i].Name < templInput.Enums[j].Name
		})

		protofile, err := templateManager.Parse("protofile", templInput)
		if err != nil {
			return err
		}

		err = WriteFile(shm.Domain, config.OutDir, protofile)
		if err != nil {
			return err
		}
	}

	return nil
}
