package parse

import (
	"errors"

	"github.com/anvlet/anvlet/cmd/schema"
)

func events(s *schema.Schema, yaml map[string]any) error {
	eventsYaml, ok := yaml["Events"]
	if !ok {
		return nil
	}

	yamlInterface, ok := eventsYaml.(map[string]interface{})
	if !ok {
		return errors.New("fail to parse Events")
	}

	schemaEvents := schema.Events{}
	for k, v := range yamlInterface {
		vMap := v.(map[string]any)

		var formats []string = nil
		if val, ok := vMap["Formats"]; ok {
			valSlice := val.([]any)
			for _, v := range valSlice {
				formats = append(formats, v.(string))
			}
		}

		fields, err := resolveField(s, vMap["Fields"].(map[string]any))
		if err != nil {
			return err
		}

		schemaEvents[k] = &schema.Event{
			Formats: formats,
			Fields:  fields,
		}
	}

	s.Events = &schemaEvents

	return nil
}
