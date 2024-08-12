package parse

import (
	"errors"
	"strconv"

	"github.com/anvil/anvil/cmd/schema"
)

func enums(s *schema.Schema, yaml map[string]any) error {
	enumsYaml, ok := yaml["Enums"]
	if !ok {
		return nil
	}

	yamlInterface, ok := enumsYaml.(map[string]any)
	if !ok {
		return errors.New("fail to parse Enums")
	}

	schemaEnums := schema.Enums{}
	for k, v := range yamlInterface {
		if k == "$ref" {
			return errors.New("cannot use $ref in Enums")
		}

		vMap, ok := v.(map[string]any)
		if !ok {
			return errors.New("fail to parse Enums." + k)
		}

		var fieldType string
		if val, ok := vMap["Type"]; ok {
			fieldType = val.(string)
		}

		var values []*schema.EnumValue = nil
		if val, ok := vMap["Values"]; ok {
			valSlice := val.([]any)
			for kk, vv := range valSlice {
				vvMap, ok := vv.(map[string]any)
				if !ok {
					return errors.New("fail to parse Enums." + k + ".Values." + strconv.Itoa(kk))
				}

				var name string
				if val, ok := vvMap["Name"]; ok {
					name = val.(string)
				}

				var value string
				if val, ok := vvMap["Value"]; ok {
					value = val.(string)
				}

				values = append(values, &schema.EnumValue{
					Name:  name,
					Value: value,
				})
			}
		}

		schemaEnums[k] = &schema.Enum{
			Type:   fieldType,
			Values: values,
		}
	}

	s.Enums = &schemaEnums

	return nil
}
