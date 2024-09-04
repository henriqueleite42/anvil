package parser

import (
	"fmt"
	"slices"

	"github.com/henriqueleite42/anvil/cli/internal/hashing"
	"github.com/henriqueleite42/anvil/cli/schemas"
)

func (self *anvToAnvpParser) usecase(file map[string]any) error {
	path := "Usecase"

	usecaseAny, ok := file["Usecase"]
	if !ok {
		return nil
	}
	usecaseMap, ok := usecaseAny.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `map[string]any`", path)
	}

	// TODO dependencies
	// TODO inputs

	methodsAny, ok := usecaseMap["Methods"]
	if !ok {
		return fmt.Errorf("\"Methods\" is a required property to \"%s\"", path)
	}
	methodsMap, ok := methodsAny.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s.Methods\" to `map[string]any`", path)
	}

	methods := &schemas.UsecaseMethods{
		Methods: map[string]*schemas.UsecaseMethod{},
	}

	// Necessary to keep some kind of order
	keys := make([]string, 0, len(methodsMap))
	for key := range methodsMap {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	for _, k := range keys {
		v := methodsMap[k]
		vMap, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("fail to parse \"%s.Methods.%s\" to `map[string]any`", path, k)
		}

		var description *string = nil
		descriptionAny, ok := vMap["Description"]
		if ok {
			descriptionString, ok := descriptionAny.(string)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.Methods.%s.Description\" to `map[string]any`", path, k)
			}
			description = &descriptionString
		}

		var input *schemas.UsecaseMethodInput = nil
		inputAny, ok := vMap["Input"]
		if ok {
			inputMap, ok := inputAny.(map[string]any)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.Methods.%s.Input\" to `map[string]any`", path, k)
			}

			var typeHash string
			typeHash, err := self.resolveType(&resolveInput{
				path: fmt.Sprintf("%s.Methods.%s", path, k),
				ref:  self.getRef("Usecase", k),
				k:    k + "Input",
				v:    inputMap,
			})
			if err != nil {
				return err
			}

			input = &schemas.UsecaseMethodInput{
				TypeHash: typeHash,
			}
		}

		var output *schemas.UsecaseMethodOutput = nil
		outputAny, ok := vMap["Output"]
		if ok {
			outputMap, ok := outputAny.(map[string]any)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.Methods.%s.Output\" to `map[string]any`", path, k)
			}

			var typeHash string
			typeHash, err := self.resolveType(&resolveInput{
				path: fmt.Sprintf("%s.Methods.%s", path, k),
				ref:  self.getRef("Usecase", k),
				k:    k + "Output",
				v:    outputMap,
			})
			if err != nil {
				return err
			}

			output = &schemas.UsecaseMethodOutput{
				TypeHash: typeHash,
			}
		}

		// TODO implement EventHashes

		ref := self.getRef(path, k)
		fullPath := fmt.Sprintf("%s.Methods.Methods.%s", path, k)

		order := len(methods.Methods)
		method := &schemas.UsecaseMethod{
			Ref:          ref,
			OriginalPath: fullPath,
			Order:        order,
			Name:         k,
			Description:  description,
			Input:        input,
			Output:       output,
		}

		methodStateHash, err := hashing.Struct(method)
		if err != nil {
			return fmt.Errorf("fail to get state hash for \"%s.Methods.%s\"", path, k)
		}
		method.StateHash = methodStateHash

		refHash := hashing.String(ref)
		methods.Methods[refHash] = method
	}

	methodsStateHash, err := hashing.Struct(methods)
	if err != nil {
		return fmt.Errorf("fail to get state hash for \"%s\"", path)
	}
	methods.StateHash = methodsStateHash

	usecase := &schemas.Usecase{
		Methods: methods,
	}

	usecaseStateHash, err := hashing.Struct(usecase)
	if err != nil {
		return fmt.Errorf("fail to get state hash for \"%s\"", path)
	}
	usecase.StateHash = usecaseStateHash

	self.schema.Usecase = usecase

	return nil
}
