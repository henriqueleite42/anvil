package internal

import (
	"fmt"
	"sort"
	"strings"

	"github.com/henriqueleite42/anvil/generators/grpc/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/grpc"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	"github.com/henriqueleite42/anvil/language-helpers/golang/template"
	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

type parser struct {
	schema                  *schemas.AnvpSchema
	grpcParser              grpc.GrpcParser
	imports                 map[string]bool
	methods                 []*templates.ProtofileTemplInputMethod
	enums                   map[string]*templates.ProtofileTemplInputEnum
	typesToAvoidDuplication map[string]*templates.ProtofileTemplInputType
	types                   []*templates.ProtofileTemplInputType
}

var templatesNamesValues = map[string]string{
	"protofile": templates.ProtofileTempl,
}

func Parse(schema *schemas.AnvpSchema, silent bool, outputFolderPath string) error {
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

	for curDomain, v := range schema.Deliveries.Deliveries {
		if _, ok := schema.Usecases.Usecases[curDomain]; !ok {
			return fmt.Errorf("no usecases specified in domain \"%s\"", curDomain)
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

		goTypeParser, err := types_parser.NewTypeParser(&types_parser.NewTypeParserInput{
			Schema:        schema,
			EnumsPkg:      "pb",
			TypesPkg:      "pb",
			EventsPkg:     "pb",
			EntitiesPkg:   "pb",
			RepositoryPkg: "pb",
			UsecasePkg:    "pb",
		})
		if err != nil {
			return err
		}

		grpcParser := grpc.NewGrpcParser(&grpc.NewGrpcParserInput{
			Schema:       schema,
			GoTypeParser: goTypeParser,
		})

		parserInstance := &parser{
			schema:                  schema,
			grpcParser:              grpcParser,
			imports:                 map[string]bool{},
			methods:                 make([]*templates.ProtofileTemplInputMethod, 0, len(rpcs)),
			enums:                   map[string]*templates.ProtofileTemplInputEnum{},
			typesToAvoidDuplication: map[string]*templates.ProtofileTemplInputType{},
			types:                   []*templates.ProtofileTemplInputType{},
		}

		emptyMsg := "google.protobuf.Empty"

		if schema.Enums != nil && schema.Enums.Enums != nil {
			for _, v := range schema.Enums.Enums {
				if strings.HasPrefix(v.Ref, curDomain) {
					parserInstance.resolveEnum(v)
				}
			}
		}

		if schema.Types != nil && schema.Types.Types != nil {
			for _, v := range schema.Types.Types {
				if strings.HasPrefix(v.Ref, curDomain) {
					parserInstance.resolveType(curDomain, v)
				}
			}
		}

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

				inputTypeResolved, err := parserInstance.resolveType(curDomain, inputType)
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

				outputTypeResolved, err := parserInstance.resolveType(curDomain, outputType)
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
			Domain:  curDomain,
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

		err = WriteFile(curDomain, outputFolderPath, protofile)
		if err != nil {
			return err
		}
	}

	return nil
}
