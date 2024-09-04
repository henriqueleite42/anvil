package parser

import (
	"fmt"

	"github.com/henriqueleite42/anvil/cli/hashing"
	"github.com/henriqueleite42/anvil/cli/schemas"
)

func (self *anvToAnvpParser) delivery(file map[string]any) error {
	path := "Delivery"

	deliveryAny, ok := file["Delivery"]
	if !ok {
		return nil
	}

	deliveryMap, ok := deliveryAny.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `map[string]any`", path)
	}

	if self.schema.Delivery == nil {
		self.schema.Delivery = &schemas.Delivery{}
	}

	grpcAny, ok := deliveryMap["Grpc"]
	if ok {
		grpcMap, ok := grpcAny.(map[string]any)
		if !ok {
			return fmt.Errorf("fail to parse \"%s.Grpc\" to `map[string]any`", path)
		}

		if self.schema.Delivery.Grpc == nil {
			self.schema.Delivery.Grpc = &schemas.DeliveryGrpc{}
		}
		if self.schema.Delivery.Grpc.Rpcs == nil {
			self.schema.Delivery.Grpc.Rpcs = map[string]*schemas.DeliveryGrpcRpc{}
		}

		rpcsAny, ok := grpcMap["Rpcs"]
		if !ok {
			return fmt.Errorf("\"Rpcs\" is a required property to \"%s.Grpc\"", path)
		}
		rpcsArr, ok := rpcsAny.([]any)
		if !ok {
			return fmt.Errorf("fail to parse \"%s.Grpc.Rpcs\" to `[]any`", path)
		}

		for k, v := range rpcsArr {
			vMap, ok := v.(map[string]any)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.Grpc.Rpcs.%d\" to `map[string]any`", path, k)
			}

			usecaseMethodAny, ok := vMap["UsecaseMethod"]
			if !ok {
				return fmt.Errorf("\"UsecaseMethod\" is a required property to \"%s.Grpc.Rpcs.%d\"", path, k)
			}
			usecaseMethodString, ok := usecaseMethodAny.(string)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.Grpc.Rpcs.%d\" to `string`", path, k)
			}

			// TODO parse examples

			// Deliveries don't use ref or paths, since they should be absolute
			// and not used by relationships

			usecaseMethod := self.getRef("Usecase", usecaseMethodString)
			usecaseMethodHash := hashing.String(usecaseMethod)

			originalPath := fmt.Sprintf("%s.Grpc.Rpcs.%d", path, k)
			ref := hashing.String(originalPath)

			rpc := &schemas.DeliveryGrpcRpc{
				Ref:               ref,
				OriginalPath:      originalPath,
				UsecaseMethodHash: usecaseMethodHash,
				Order:             k,
			}

			stateHash, err := hashing.Struct(rpc)
			if err != nil {
				return fmt.Errorf("fail to get state hash for \"%s.Grpc.Rpcs.%d\"", path, k)
			}
			rpc.StateHash = stateHash

			self.schema.Delivery.Grpc.Rpcs[ref] = rpc
		}

		stateHash, err := hashing.Struct(self.schema.Delivery.Grpc)
		if err != nil {
			return fmt.Errorf("fail to get state hash for \"%s.Grpc\"", path)
		}
		self.schema.Delivery.Grpc.StateHash = stateHash
	}

	// TODO parse http

	// TODO parse queue

	stateHash, err := hashing.Struct(self.schema.Delivery)
	if err != nil {
		return fmt.Errorf("fail to get state hash for \"%s\"", path)
	}
	self.schema.Delivery.StateHash = stateHash

	return nil
}
