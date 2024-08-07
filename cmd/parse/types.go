package parse

import (
	"errors"

	"github.com/anuntech/hephaestus/cmd/schema"
)

func types(s *schema.Schema, yaml map[string]any) error {
	typesYaml, ok := yaml["Types"]
	if !ok {
		return nil
	}

	yamlInterface, ok := typesYaml.(map[string]interface{})
	if !ok {
		return errors.New("fail to parse Types")
	}

	schemaTypes := schema.Types{}
	for k, v := range yamlInterface {
		if k == "$ref" {
			return errors.New("cannot use $ref in Types")
		}

		fieldInterface, ok := v.(map[string]interface{})
		if !ok {
			return errors.New("fail to parse Types." + k)
		}

		typesResolved, err := resolveField(s, fieldInterface)
		if err != nil {
			return err
		}

		schemaTypes[k] = typesResolved
	}

	s.Types = &schemaTypes

	return nil
}
