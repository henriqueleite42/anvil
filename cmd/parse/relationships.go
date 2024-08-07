package parse

import (
	"errors"

	"github.com/anuntech/hephaestus/cmd/schema"
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
		relationshipVal, ok := v.(string)
		if !ok {
			return errors.New("fail to parse Relationships." + k)
		}
		schemaRelationships[k] = relationshipVal
	}

	s.Relationships = &schemaRelationships

	return nil
}
