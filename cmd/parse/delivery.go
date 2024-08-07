package parse

import (
	"errors"

	"github.com/anuntech/hephaestus/cmd/schema"
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

	var grpc *schema.DeliveryGrpc = nil
	grpcMap, ok := yamlInterface["Grpc"].(map[string]any)
	if ok {
		var genClient *bool = nil
		if val, ok := grpcMap["GenClient"]; ok {
			valBool := val.(bool)
			genClient = &valBool
		}
		var clientPath *string = nil
		if val, ok := grpcMap["ClientPath"]; ok {
			valString := val.(string)
			clientPath = &valString
		}
		var genProto *bool = nil
		if val, ok := grpcMap["GenProto"]; ok {
			valBool := val.(bool)
			genProto = &valBool
		}
		var protofilePath *string = nil
		if val, ok := grpcMap["ProtofilePath"]; ok {
			valString := val.(string)
			protofilePath = &valString
		}
		var genServer *string = nil
		if val, ok := grpcMap["GenServer"]; ok {
			valString := val.(string)
			genServer = &valString
		}
		var serverPath *string = nil
		if val, ok := grpcMap["ServerPath"]; ok {
			valString := val.(string)
			serverPath = &valString
		}

		grpc = &schema.DeliveryGrpc{
			GenClient:     genClient,
			ClientPath:    clientPath,
			GenProto:      genProto,
			ProtofilePath: protofilePath,
			GenServer:     genServer,
			ServerPath:    serverPath,
		}
	}

	s.Delivery = &schema.Delivery{
		Dependencies: dependencies,
		Grpc:         grpc,
	}

	return nil
}
