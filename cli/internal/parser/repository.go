package parser

import (
	"fmt"
	"sort"

	"github.com/henriqueleite42/anvil/cli/hashing"
	"github.com/henriqueleite42/anvil/cli/schemas"
)

func (self *anvToAnvpParser) repository(file map[string]any) error {
	path := "Repository"

	repositoryAny, ok := file["Repository"]
	if !ok {
		return nil
	}
	repositoryMap, ok := repositoryAny.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `map[string]any`", path)
	}

	// TODO dependencies
	// TODO inputs

	methodsAny, ok := repositoryMap["Methods"]
	if !ok {
		return fmt.Errorf("\"Methods\" is a required property to \"%s\"", path)
	}
	methodsMap, ok := methodsAny.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s.Methods\" to `map[string]any`", path)
	}

	methods := &schemas.RepositoryMethods{
		Methods: map[string]*schemas.RepositoryMethod{},
	}

	// Necessary to keep some kind of order
	keys := []string{}
	for key := range methodsMap {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

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

		var input *schemas.RepositoryMethodInput = nil
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

			input = &schemas.RepositoryMethodInput{
				TypeHash: typeHash,
			}
		}

		var output *schemas.RepositoryMethodOutput = nil
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

			output = &schemas.RepositoryMethodOutput{
				TypeHash: typeHash,
			}
		}

		fullPath := fmt.Sprintf("%s.Methods.Methods.%s", path, k)

		order := len(methods.Methods)
		method := &schemas.RepositoryMethod{
			Ref:          self.getRef("", fmt.Sprintf("%s.%s", path, k)),
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

		methodHash := hashing.String(fullPath)
		methods.Methods[methodHash] = method
	}

	methodsStateHash, err := hashing.Struct(methods)
	if err != nil {
		return fmt.Errorf("fail to get state hash for \"%s\"", path)
	}
	methods.StateHash = methodsStateHash

	repository := &schemas.Repository{
		Methods: methods,
	}

	repositoryStateHash, err := hashing.Struct(repository)
	if err != nil {
		return fmt.Errorf("fail to get state hash for \"%s\"", path)
	}
	repository.StateHash = repositoryStateHash

	self.schema.Repository = repository

	return nil
}
