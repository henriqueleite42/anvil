package parse

import (
	"errors"

	"github.com/anvlet/anvlet/cmd/schema"
)

func repository(s *schema.Schema, yaml map[string]any) error {
	repositoryYaml, ok := yaml["Repository"]
	if !ok {
		return nil
	}

	yamlInterface, ok := repositoryYaml.(map[string]interface{})
	if !ok {
		return errors.New("fail to parse Repository")
	}

	var dependencies map[string]*schema.Dependency = nil
	dependenciesMap, ok := yamlInterface["Dependencies"].(map[string]any)
	if ok {
		dependencies = map[string]*schema.Dependency{}
		for k, v := range dependenciesMap {
			vMap := v.(map[string]any)
			dependency, err := parseDependency(vMap)
			if err != nil {
				return errors.New("fail to parse Repository Dependency")
			}

			dependencies[k] = dependency
		}
	}

	var inputs map[string]*schema.Dependency = nil
	inputsMap, ok := yamlInterface["Inputs"].(map[string]any)
	if ok {
		inputs = map[string]*schema.Dependency{}
		for k, v := range inputsMap {
			vMap := v.(map[string]any)
			dependency, err := parseDependency(vMap)
			if err != nil {
				return errors.New("fail to parse Repository Input")
			}

			inputs[k] = dependency
		}
	}

	methods := map[string]*schema.RepositoryMethod{}
	methodsMap, _ := yamlInterface["Methods"].(map[string]any)
	for k, v := range methodsMap {
		vMap := v.(map[string]any)

		var inputs map[string]*schema.Field = nil
		if _, ok := vMap["Input"]; ok {
			inputsMap := vMap["Input"].(map[string]any)
			fields, err := resolveField(s, inputsMap)
			if err != nil {
				return errors.New("fail to parse Repository method input")
			}
			inputs = fields
		}

		var outputs map[string]*schema.Field = nil
		if _, ok := vMap["Output"]; ok {
			outputsMap := vMap["Output"].(map[string]any)
			fields, err := resolveField(s, outputsMap)
			if err != nil {
				return errors.New("fail to parse Repository method output")
			}
			outputs = fields
		}

		methods[k] = &schema.RepositoryMethod{
			Input:  inputs,
			Output: outputs,
		}
	}

	s.Repository = &schema.Repository{
		Dependencies: dependencies,
		Inputs:       inputs,
		Methods:      methods,
	}

	return nil
}
