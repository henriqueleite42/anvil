package parser_anv

import (
	"fmt"

	"github.com/anvil/anvil/internal/hashing"
	"github.com/anvil/anvil/schemas"
)

func (self *anvToAnvpParser) usecase(file map[string]any) error {
	path := self.getPath("Usecase")

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

	for k, v := range methodsMap {
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
				ref:  self.getRef("", fmt.Sprintf("%s.%s", path, k)),
				k:    "Input",
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
				ref:  self.getRef("", fmt.Sprintf("%s.Methods.%s.%s", path, k, "Output")),
				k:    "Output",
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

		fullPath := self.getPath(fmt.Sprintf("%s.Methods.Methods.%s", path, k))

		method := &schemas.UsecaseMethod{
			Ref:          self.getRef("", fmt.Sprintf("%s.%s", path, k)),
			OriginalPath: fullPath,
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

		methodHash := hashing.String(fullPath)
		methods.Methods[methodHash] = method
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
