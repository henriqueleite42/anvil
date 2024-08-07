package parse

import (
	"errors"

	"github.com/anuntech/hephaestus/cmd/schema"
)

func usecase(s *schema.Schema, yaml map[string]any) error {
	usecaseYaml, ok := yaml["Usecase"]
	if !ok {
		return nil
	}

	yamlInterface, ok := usecaseYaml.(map[string]interface{})
	if !ok {
		return errors.New("fail to parse Usecase")
	}

	dependencies := map[string]*schema.Dependency{}
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

	methods := map[string]*schema.Method{}
	methodsMap, _ := yamlInterface["Methods"].(map[string]any)
	for k, v := range methodsMap {
		vMap := v.(map[string]any)

		var inputs map[string]*schema.Field = nil
		if _, ok := vMap["Input"]; ok {
			inputsMap := vMap["Input"].(map[string]any)
			fields, err := resolveField(s, inputsMap)
			if err != nil {
				return errors.New("fail to parse Usecase method input")
			}
			inputs = fields
		}

		var outputs map[string]*schema.Field = nil
		if _, ok := vMap["Output"]; ok {
			outputsMap := vMap["Output"].(map[string]any)
			fields, err := resolveField(s, outputsMap)
			if err != nil {
				return errors.New("fail to parse Usecase method output")
			}
			outputs = fields
		}

		var delivery *schema.MethodDelivery = nil
		if _, ok := vMap["Delivery"]; ok {
			deliveryMap := vMap["Delivery"].(map[string]any)

			delivery = &schema.MethodDelivery{}
			if grpcAny, ok := deliveryMap["Grpc"]; ok {
				var client *bool = nil
				examples := map[string]*schema.MethodDeliveryGrpc_Example{}

				if grpcMap, ok := grpcAny.(map[string]any); ok {
					if clientAny, ok := grpcMap["Client"]; ok {
						clientBool := clientAny.(bool)
						client = &clientBool
					}

					if examplesAny, ok := grpcMap["Examples"]; ok {
						examplesMap := examplesAny.(map[string]any)
						for k, v := range examplesMap {
							vMap := v.(map[string]any)

							var statusCode int
							if val, ok := vMap["StatusCode"]; ok {
								statusCode = val.(int)
							}
							var message any
							if val, ok := vMap["Message"]; ok {
								message = val
							}
							var returns any
							if val, ok := vMap["Returns"]; ok {
								returns = val
							}

							examples[k] = &schema.MethodDeliveryGrpc_Example{
								StatusCode: statusCode,
								Message:    message,
								Returns:    returns,
							}
						}
					}
				}

				delivery.Grpc = &schema.MethodDeliveryGrpc{
					Client:   client,
					Examples: examples,
				}
			}

			if queueAny, ok := deliveryMap["Queue"]; ok {
				var id string
				var relatedTo string

				if queueMap, ok := queueAny.(map[string]any); ok {
					if idAny, ok := queueMap["Id"]; ok {
						idString := idAny.(string)
						id = idString
					}

					if relatedToAny, ok := queueMap["RelatedTo"]; ok {
						relatedToString := relatedToAny.(string)
						relatedTo = relatedToString
					}
				}

				delivery.Queue = &schema.MethodDeliveryQueue{
					Id:        id,
					RelatedTo: relatedTo,
				}
			}
		}

		methods[k] = &schema.Method{
			Input:    inputs,
			Output:   outputs,
			Delivery: delivery,
		}
	}

	s.Usecase = &schema.Usecase{
		Dependencies: dependencies,
		Methods:      methods,
	}

	return nil
}
