package parse

import (
	"errors"

	"github.com/anuntech/hephaestus/cmd/types"
)

func Events(s *types.Schema, yaml map[string]any) error {
	eventsYaml, ok := yaml["Events"]
	if !ok {
		return nil
	}

	yamlInterface, ok := eventsYaml.(map[string]interface{})
	if !ok {
		return errors.New("fail to parse Events")
	}

	schemaEvents := types.Events{}
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

		schemaEvents[k] = &types.Event{
			Formats: formats,
			Fields:  fields,
		}
	}

	s.Events = &schemaEvents

	return nil
}
