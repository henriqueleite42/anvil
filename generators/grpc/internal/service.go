package internal

import (
	"fmt"
	"sort"
	"strings"
)

const EMPTY_TYPE = "google.protobuf.Empty"

func (self *protoFile) resolveService() error {
	if self.schema.Domain == "" {
		return fmt.Errorf("no domain specified")
	}
	if self.schema.Delivery == nil {
		return fmt.Errorf("no delivery specified")
	}
	if self.schema.Delivery.Grpc == nil {
		return fmt.Errorf("no gRPC delivery specified")
	}
	if self.schema.Delivery.Grpc.Rpcs == nil {
		return fmt.Errorf("no RPCs specified for gRPC delivery")
	}
	if self.schema.Usecase == nil {
		return fmt.Errorf("no usecases to deliver")
	}
	if self.schema.Usecase.Methods == nil || self.schema.Usecase.Methods.Methods == nil {
		return fmt.Errorf("no usecases methods to deliver")
	}

	sortedRpcs := []*SortedByOrder{}
	for k, v := range self.schema.Delivery.Grpc.Rpcs {
		sortedRpcs = append(sortedRpcs, &SortedByOrder{
			Order: v.Order,
			Key:   k,
		})
	}
	sort.Slice(sortedRpcs, func(i, j int) bool {
		return sortedRpcs[i].Order < sortedRpcs[j].Order
	})

	methods := []string{}
	for _, sortedRpc := range sortedRpcs {
		k := sortedRpc.Key
		v := self.schema.Delivery.Grpc.Rpcs[k]

		if v.UsecaseMethodHash == "" {
			return fmt.Errorf("missing \"UsecaseMethodHash\" for RPC \"%s\"", k)
		}

		usecase, ok := self.schema.Usecase.Methods.Methods[v.UsecaseMethodHash]
		if !ok {
			return fmt.Errorf("usecase method \"%s\" not found", v.UsecaseMethodHash)
		}

		input := EMPTY_TYPE
		if usecase.Input != nil && usecase.Input.TypeHash != "" {
			uscType, err := self.resolveMsgPropType(nil, usecase.Input.TypeHash)
			if err != nil {
				return err
			}
			input = uscType
		}

		output := EMPTY_TYPE
		if usecase.Output != nil && usecase.Output.TypeHash != "" {
			uscType, err := self.resolveMsgPropType(nil, usecase.Output.TypeHash)
			if err != nil {
				return err
			}
			output = uscType
		}

		if input == EMPTY_TYPE || output == EMPTY_TYPE {
			self.imports["import \"google/protobuf/empty.proto\";"] = true
		}

		methods = append(methods, fmt.Sprintf("	rpc %s(%s) returns (%s) {}", usecase.Name, input, output))
	}

	serviceString := fmt.Sprintf(`service %s {
%s
}`, self.schema.Domain, strings.Join(methods, "\n"))

	self.service = serviceString

	return nil
}
