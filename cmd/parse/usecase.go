package parse

import (
	"errors"

	"github.com/anuntech/hephaestus/cmd/types"
)

func Usecase(s *types.Schema, yaml map[string]any) error {
	usecaseYaml, ok := yaml["Usecase"]
	if !ok {
		return nil
	}

	yamlInterface, ok := usecaseYaml.(map[string]interface{})
	if !ok {
		return errors.New("fail to parse Usecase")
	}

	dependencies := map[string]*types.Dependency{}
	dependenciesMap, ok := yamlInterface["Dependencies"].(map[string]any)
	if ok {
		for k, v := range dependenciesMap {
			vMap := v.(map[string]any)
			dependency, err := parseDependency(vMap)
			if err != nil {
				return errors.New("fail to parse Usecase Dependency")
			}

			dependencies[k] = dependency
		}
	}

	methods := map[string]*types.Method{}
	methodsMap, _ := yamlInterface["Methods"].(map[string]any)
	for k, v := range methodsMap {
		vMap := v.(map[string]any)

		var inputs map[string]*types.Field = nil
		if _, ok := vMap["Input"]; ok {
			inputsMap := vMap["Input"].(map[string]any)
			fields, err := resolveField(s, inputsMap)
			if err != nil {
				return errors.New("fail to parse Usecase method input")
			}
			inputs = fields
		}

		var outputs map[string]*types.Field = nil
		if _, ok := vMap["Output"]; ok {
			outputsMap := vMap["Output"].(map[string]any)
			fields, err := resolveField(s, outputsMap)
			if err != nil {
				return errors.New("fail to parse Usecase method output")
			}
			outputs = fields
		}

		var delivery *types.MethodDelivery = nil
		if _, ok := vMap["Delivery"]; ok {
			deliveryMap := vMap["Delivery"].(map[string]any)

			delivery = &types.MethodDelivery{}
			if grpcAny, ok := deliveryMap["Grpc"]; ok {
				var client *bool = nil

				if grpcMap, ok := grpcAny.(map[string]any); ok {
					if clientAny, ok := grpcMap["Client"]; ok {
						clientBool := clientAny.(bool)
						client = &clientBool
					}
				}

				delivery.Grpc = &types.MethodDeliveryGrpc{
					Client: client,
				}
			}
		}

		methods[k] = &types.Method{
			Input:    inputs,
			Output:   outputs,
			Delivery: delivery,
		}
	}

	s.Usecase = &types.Usecase{
		Dependencies: dependencies,
		Methods:      methods,
	}

	return nil
}
