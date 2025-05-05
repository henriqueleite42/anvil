package parser

import (
	"fmt"

	"github.com/ettle/strcase"
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

	httpAny, ok := deliveryMap["Http"]
	if ok {
		httpMap, ok := httpAny.(map[string]any)
		if !ok {
			return fmt.Errorf("fail to parse \"%s.Http\" to `map[string]any`", path)
		}

		self.schema.Deliveries.Deliveries[curDomain].Http = &schemas.DeliveryHttp{
			Routes: map[string]*schemas.DeliveryHttpRoute{},
		}

		routesAny, ok := httpMap["Routes"]
		if !ok {
			return fmt.Errorf("\"Queues\" is a required property to \"%s.Queue\"", path)
		}
		routesArr, ok := routesAny.([]any)
		if !ok {
			return fmt.Errorf("fail to parse \"%s.Http.Routes\" to `[]any`", path)
		}

		for k, v := range routesArr {
			vMap, ok := v.(map[string]any)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.Http.Routes.%d\" to `map[string]any`", path, k)
			}

			usecaseMethodAny, ok := vMap["UsecaseMethod"]
			if !ok {
				return fmt.Errorf("\"UsecaseMethod\" is a required property to \"%s.Http.Routes.%d\"", path, k)
			}
			usecaseMethodString, ok := usecaseMethodAny.(string)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.Http.Routes.%d.UsecaseMethod\" to `string`", path, k)
			}

			usecaseMethodRef := self.getRef(curDomain, "Usecase."+usecaseMethodString)
			usecaseMethodHash := hashing.String(usecaseMethodRef)

			pathAny, ok := vMap["Path"]
			if !ok {
				return fmt.Errorf("\"Path\" is a required property to \"%s.Http.Routes.%d\"", path, k)
			}
			pathString, ok := pathAny.(string)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.Http.Routes.%d.Path\" to `string`", path, k)
			}

			methodAny, ok := vMap["Method"]
			if !ok {
				return fmt.Errorf("\"Method\" is a required property to \"%s.Http.Routes.%d\"", path, k)
			}
			methodString, ok := methodAny.(string)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.Http.Routes.%d.Method\" to `string`", path, k)
			}

			// TODO parse other properties

			originalPath := fmt.Sprintf("%s.QHttp.Routes.%d", path, k)
			ref := self.getRef(curDomain, "Http.Routes."+strcase.ToPascal(pathString))

			route := &schemas.DeliveryHttpRoute{
				Ref:               ref,
				OriginalPath:      originalPath,
				Domain:            curDomain,
				UsecaseMethodHash: usecaseMethodHash,
				Path:              pathString,
				HttpMethod:        methodString,
			}

			stateHash, err := hashing.Struct(route)
			if err != nil {
				return fmt.Errorf("fail to get state hash for \"%s.Http.Routes.%d\"", path, k)
			}
			route.StateHash = stateHash

			self.schema.Deliveries.Deliveries[curDomain].Http.Routes[ref] = route
		}

		// Validate duplicated Routes
		existentRoutes := make(map[string]bool, len(self.schema.Deliveries.Deliveries[curDomain].Http.Routes))
		for _, v := range self.schema.Deliveries.Deliveries[curDomain].Http.Routes {
			if existentRoutes[v.Ref] {
				return fmt.Errorf("duplicated http route \"%s\"", v.OriginalPath)
			}

			existentRoutes[v.Ref] = true
		}

		stateHash, err := hashing.Struct(self.schema.Deliveries.Deliveries[curDomain].Http)
		if err != nil {
			return fmt.Errorf("fail to get state hash for \"%s.Http\"", path)
		}
		self.schema.Deliveries.Deliveries[curDomain].Http.StateHash = stateHash
	}

	queueAny, ok := deliveryMap["Queue"]
	if ok {
		queueMap, ok := queueAny.(map[string]any)
		if !ok {
			return fmt.Errorf("fail to parse \"%s.Queue\" to `map[string]any`", path)
		}

		self.schema.Deliveries.Deliveries[curDomain].Queue = &schemas.DeliveryQueue{
			Queues: map[string]*schemas.DeliveryQueueQueue{},
		}

		queuesAny, ok := queueMap["Queues"]
		if !ok {
			return fmt.Errorf("\"Queues\" is a required property to \"%s.Queue\"", path)
		}
		queuesArr, ok := queuesAny.([]any)
		if !ok {
			return fmt.Errorf("fail to parse \"%s.Queue.Queues\" to `[]any`", path)
		}

		for k, v := range queuesArr {
			vMap, ok := v.(map[string]any)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.Queue.Queues.%d\" to `map[string]any`", path, k)
			}

			usecaseMethodAny, ok := vMap["UsecaseMethod"]
			if !ok {
				return fmt.Errorf("\"UsecaseMethod\" is a required property to \"%s.Queue.Queues.%d\"", path, k)
			}
			usecaseMethodString, ok := usecaseMethodAny.(string)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.Queue.Queues.%d.UsecaseMethod\" to `string`", path, k)
			}

			usecaseMethodRef := self.getRef(curDomain, "Usecase."+usecaseMethodString)
			usecaseMethodHash := hashing.String(usecaseMethodRef)

			idAny, ok := vMap["Id"]
			if !ok {
				return fmt.Errorf("\"Id\" is a required property to \"%s.Queue.Queues.%d\"", path, k)
			}
			idString, ok := idAny.(string)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.Queue.Queues.%d.Id\" to `string`", path, k)
			}

			var bulk bool
			bulkAny, ok := vMap["Bulk"]
			if ok {
				bulkBool, ok := bulkAny.(bool)
				if !ok {
					return fmt.Errorf("fail to parse \"%s.Queue.Queues.%d.Bulk\" to `bool`", path, k)
				}
				bulk = bulkBool
			}

			// TODO parse examples

			originalPath := fmt.Sprintf("%s.Queue.Queues.%d", path, k)
			ref := self.getRef(curDomain, "Queue.Queues."+strcase.ToPascal(idString))

			queue := &schemas.DeliveryQueueQueue{
				Ref:               ref,
				OriginalPath:      originalPath,
				Domain:            curDomain,
				UsecaseMethodHash: usecaseMethodHash,
				QueueId:           idString,
				Bulk:              bulk,
			}

			stateHash, err := hashing.Struct(queue)
			if err != nil {
				return fmt.Errorf("fail to get state hash for \"%s.Queue.Queues.%d\"", path, k)
			}
			queue.StateHash = stateHash

			self.schema.Deliveries.Deliveries[curDomain].Queue.Queues[ref] = queue
		}

		// Validate duplicated Queues
		existentQueues := make(map[string]bool, len(self.schema.Deliveries.Deliveries[curDomain].Queue.Queues))
		for _, v := range self.schema.Deliveries.Deliveries[curDomain].Queue.Queues {
			if existentQueues[v.Ref] {
				return fmt.Errorf("duplicated queue \"%s\"", v.OriginalPath)
			}

			existentQueues[v.Ref] = true
		}

		stateHash, err := hashing.Struct(self.schema.Deliveries.Deliveries[curDomain].Queue)
		if err != nil {
			return fmt.Errorf("fail to get state hash for \"%s.Queue\"", path)
		}
		self.schema.Deliveries.Deliveries[curDomain].Queue.StateHash = stateHash
	}

	stateHash, err := hashing.Struct(self.schema.Deliveries.Deliveries[curDomain])
	if err != nil {
		return fmt.Errorf("fail to get state hash for \"%s\"", path)
	}
	self.schema.Deliveries.Deliveries[curDomain].StateHash = stateHash

	return nil
}
