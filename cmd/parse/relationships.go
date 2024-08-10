package parse

import (
	"errors"

	"github.com/anvlet/anvlet/cmd/schema"
)

func relationships(s *schema.Schema, yaml map[string]any) error {
	relationships, ok := yaml["Relationships"]
	if !ok {
		return nil
	}

	yamlInterface, ok := relationships.(map[string]any)
	if !ok {
		return errors.New("fail to parse Relationships")
	}

	schemaRelationships := schema.Relationships{}
	for k, v := range yamlInterface {
		vMap := v.(map[string]any)

		var uri string
		if val, ok := vMap["Uri"]; ok {
			valString := val.(string)
			uri = valString
		}

		schemaRelationships[k] = &schema.Relationship{
			Uri: uri,
		}
	}

	s.Relationships = &schemaRelationships

	return nil
}
