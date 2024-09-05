package contract

import (
	"fmt"
	"sort"
)

func (self *contractFile) parseApi() error {
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

		input := ""
		if usecase.Input != nil && usecase.Input.TypeHash != "" {
			uscType, err := self.resolveType(usecase.Input.TypeHash, "")
			if err != nil {
				return err
			}
			input = "i " + uscType
		}

		output := " error"
		if usecase.Output != nil && usecase.Output.TypeHash != "" {
			uscType, err := self.resolveType(usecase.Output.TypeHash, "")
			if err != nil {
				return err
			}
			output = fmt.Sprintf(" (%s, error)", uscType)
		}

		self.methods = append(self.methods, fmt.Sprintf("	%s(%s)%s", usecase.Name, input, output))
	}

	return nil
}
