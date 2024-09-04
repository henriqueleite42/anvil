package internal

import (
	"fmt"
	"sort"
	"strings"

	"github.com/henriqueleite42/anvil/cli/schemas"
)

type SortedByOrder struct {
	Order int
	Key   string
}

const EMPTY_TYPE = "google.protobuf.Empty"

func (self *protoFile) resolveService(schema *schemas.Schema) error {
	if schema.Domain == "" {
		return fmt.Errorf("no domain specified")
	}
	if schema.Delivery == nil {
		return fmt.Errorf("no delivery specified")
	}
	if schema.Delivery.Grpc == nil {
		return fmt.Errorf("no gRPC delivery specified")
	}
	if schema.Delivery.Grpc.Rpcs == nil {
		return fmt.Errorf("no RPCs specified for gRPC delivery")
	}
	if schema.Usecase == nil {
		return fmt.Errorf("no usecases to deliver")
	}
	if schema.Usecase.Methods == nil || schema.Usecase.Methods.Methods == nil {
		return fmt.Errorf("no usecases methods to deliver")
	}

	sortedRpcs := []*SortedByOrder{}
	for k, v := range schema.Delivery.Grpc.Rpcs {
		sortedRpcs = append(sortedRpcs, &SortedByOrder{
			Order: v.Order,
			Key:   k,
		})
	}
	sort.Slice(sortedRpcs, func(i, j int) bool {
		return sortedRpcs[i].Order < sortedRpcs[j].Order
	})
	for _, k := range sortedRpcs {
		fmt.Println(k.Order)
	}

	methods := []string{}
	for _, sortedRpc := range sortedRpcs {
		k := sortedRpc.Key
		v := schema.Delivery.Grpc.Rpcs[k]

		if v.UsecaseMethodHash == "" {
			return fmt.Errorf("missing \"UsecaseMethodHash\" for RPC \"%s\"", k)
		}

		usecase, ok := schema.Usecase.Methods.Methods[v.UsecaseMethodHash]
		if !ok {
			return fmt.Errorf("usecase method \"%s\" not found", v.UsecaseMethodHash)
		}

		input := EMPTY_TYPE
		if usecase.Input != nil && usecase.Input.TypeHash != "" {
			uscType, err := self.resolveMsgPropType(schema, usecase.Input.TypeHash)
			if err != nil {
				return err
			}
			input = uscType
		}

		output := EMPTY_TYPE
		if usecase.Output != nil && usecase.Output.TypeHash != "" {
			uscType, err := self.resolveMsgPropType(schema, usecase.Output.TypeHash)
			if err != nil {
				return err
			}
			output = uscType
		}

		if input == EMPTY_TYPE || output == EMPTY_TYPE {
			self.imports["import \"google/protobuf/empty.proto\";"] = true
		}

		methods = append(methods, fmt.Sprintf("	%s(%s) returns (%s) {}", usecase.Name, input, output))
	}

	serviceString := fmt.Sprintf(`service %s {
%s
}`, schema.Domain, strings.Join(methods, "\n"))

	self.service = serviceString

	return nil
}
