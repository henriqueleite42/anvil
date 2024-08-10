package parse

import (
	"errors"

	"github.com/anvlet/anvlet/cmd/schema"
)

func delivery(s *schema.Schema, yaml map[string]any) error {
	deliveryYaml, ok := yaml["Delivery"]
	if !ok {
		return nil
	}

	yamlInterface, ok := deliveryYaml.(map[string]interface{})
	if !ok {
		return errors.New("fail to parse Delivery")
	}

	var dependencies map[string]*schema.Dependency = nil
	dependenciesMap, ok := yamlInterface["Dependencies"].(map[string]any)
	if ok {
		dependencies = map[string]*schema.Dependency{}
		for k, v := range dependenciesMap {
			vMap := v.(map[string]any)
			dependency, err := parseDependency(vMap)
			if err != nil {
				return errors.New("fail to parse Delivery Dependency")
			}

			dependencies[k] = dependency
		}
	}

	s.Delivery = &schema.Delivery{
		Dependencies: dependencies,
	}

	return nil
}
