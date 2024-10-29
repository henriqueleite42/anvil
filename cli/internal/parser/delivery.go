package parser

import (
	"fmt"

	"github.com/henriqueleite42/anvil/language-helpers/golang/hashing"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *anvToAnvpParser) delivery(curDomain string, file map[string]any) error {
	deliveryAny, ok := file["Delivery"]
	if !ok {
		return nil
	}

	if self.schema.Deliveries == nil {
		self.schema.Deliveries = &schemas.Deliveries{}
	}
	if self.schema.Deliveries.Deliveries == nil {
		self.schema.Deliveries.Deliveries = map[string]*schemas.Delivery{}
	}

	path := curDomain + ".Delivery"

	deliveryMap, ok := deliveryAny.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `map[string]any`", path)
	}

	self.schema.Deliveries.Deliveries[curDomain] = &schemas.Delivery{}

	var servers map[string]*schemas.DeliveryServers = nil
	serversAny, ok := deliveryMap["Servers"]
	if ok {
		valMap, ok := serversAny.(map[string]any)
		if !ok {
			return fmt.Errorf("fail to parse \"%s.Servers\" to `map[string]any`", path)
		}

		servers = make(map[string]*schemas.DeliveryServers, len(valMap))

		for k, v := range valMap {
			vMap, ok := v.(map[string]any)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.Servers.%s\" to `map[string]any`", path, k)
			}

			urlAny, ok := vMap["Url"]
			if !ok {
				return fmt.Errorf("\"Url\" is a required property to \"%s.Servers.%s\"", path, k)
			}
			urlString, ok := urlAny.(string)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.Servers.%s.Url\" to `string`", path, k)
			}

			servers[k] = &schemas.DeliveryServers{
				Url: urlString,
			}
		}

		self.schema.Deliveries.Deliveries[curDomain].Servers = servers
	}

	grpcAny, ok := deliveryMap["Grpc"]
	if ok {
		grpcMap, ok := grpcAny.(map[string]any)
		if !ok {
			return fmt.Errorf("fail to parse \"%s.Grpc\" to `map[string]any`", path)
		}

		self.schema.Deliveries.Deliveries[curDomain].Grpc = &schemas.DeliveryGrpc{
			Rpcs: map[string]*schemas.DeliveryGrpcRpc{},
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
				return fmt.Errorf("fail to parse \"%s.Grpc.Rpcs.%d.UsecaseMethod\" to `string`", path, k)
			}

			// TODO parse examples

			usecaseMethodRef := self.getRef(curDomain, "Usecase."+usecaseMethodString)
			usecaseMethodHash := hashing.String(usecaseMethodRef)

			var name string
			nameAny, ok := vMap["Name"]
			if ok {
				nameString, ok := nameAny.(string)
				if !ok {
					return fmt.Errorf("fail to parse \"%s.Grpc.Rpcs.%d.Name\" to `string`", path, k)
				}
				name = nameString
			} else {
				if self.schema.Usecases == nil ||
					self.schema.Usecases.Usecases == nil {
					return fmt.Errorf("fail to find usecase \"%s\" for rpc \"%s.Grpc.Rpcs.%d\": no usecases defined", usecaseMethodHash, path, k)
				}
				if _, ok := self.schema.Usecases.Usecases[curDomain]; !ok {
					return fmt.Errorf("fail to find usecase \"%s\" for rpc \"%s.Grpc.Rpcs.%d\": no usecases for domain \"%s\" defined", usecaseMethodHash, path, k, curDomain)
				}

				usecase, ok := self.schema.Usecases.Usecases[curDomain].Methods.Methods[usecaseMethodHash]
				if !ok {
					return fmt.Errorf("fail to find usecase \"%s\" for rpc \"%s.Grpc.Rpcs.%d\"", usecaseMethodHash, path, k)
				}

				name = usecase.Name
			}

			originalPath := fmt.Sprintf("%s.Grpc.Rpcs.%d", path, k)
			ref := self.getRef(curDomain, "Delivery.Grpc."+name)

			rpc := &schemas.DeliveryGrpcRpc{
				Ref:               ref,
				OriginalPath:      originalPath,
				Domain:            curDomain,
				Name:              name,
				UsecaseMethodHash: usecaseMethodHash,
				Order:             uint(k),
			}

			stateHash, err := hashing.Struct(rpc)
			if err != nil {
				return fmt.Errorf("fail to get state hash for \"%s.Grpc.Rpcs.%d\"", path, k)
			}
			rpc.StateHash = stateHash

			self.schema.Deliveries.Deliveries[curDomain].Grpc.Rpcs[ref] = rpc
		}

		// Validate duplicated Rpcs
		existentRpcs := make(map[string]bool, len(self.schema.Deliveries.Deliveries[curDomain].Grpc.Rpcs))
		for _, v := range self.schema.Deliveries.Deliveries[curDomain].Grpc.Rpcs {
			if existentRpcs[v.Ref] {
				return fmt.Errorf("duplicated grpc rpc \"%s\"", v.OriginalPath)
			}

			existentRpcs[v.Ref] = true
		}

		stateHash, err := hashing.Struct(self.schema.Deliveries.Deliveries[curDomain].Grpc)
		if err != nil {
			return fmt.Errorf("fail to get state hash for \"%s.Grpc\"", path)
		}
		self.schema.Deliveries.Deliveries[curDomain].Grpc.StateHash = stateHash
	}

	// TODO parse http

	// TODO parse queue

	stateHash, err := hashing.Struct(self.schema.Deliveries.Deliveries[curDomain])
	if err != nil {
		return fmt.Errorf("fail to get state hash for \"%s\"", path)
	}
	self.schema.Deliveries.Deliveries[curDomain].StateHash = stateHash

	return nil
}
