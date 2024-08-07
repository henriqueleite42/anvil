package parse

import (
	"errors"

	"github.com/anuntech/hephaestus/cmd/schema"
)

func enums(s *schema.Schema, yaml map[string]any) error {
	enums, ok := yaml["Enums"]
	if !ok {
		return nil
	}

	yamlInterface, ok := enums.(map[string]any)
	if !ok {
		return errors.New("fail to parse Enums")
	}

	schemaEnums := schema.Enums{}
	for k, v := range yamlInterface {
		enum := map[string]string{}

		enumParsed, ok := v.(map[string]any)
		if !ok {
			return errors.New("fail to parse Enums." + k)
		}

		for kk, vv := range enumParsed {
			enumVal, ok := vv.(string)
			if !ok {
				return errors.New("fail to parse Enums." + k + "." + kk)
			}
			enum[kk] = enumVal
		}

		schemaEnums[k] = enum
	}

	s.Enums = &schemaEnums

	return nil
}
