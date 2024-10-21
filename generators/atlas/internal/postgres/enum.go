package postgres

import (
	"sort"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func resolveEnums(schema *schemas.AnvpSchema) ([]*schemas.Enum, error) {
	if schema.Enums == nil || schema.Enums.Enums == nil {
		return []*schemas.Enum{}, nil
	}

	enums := make([]*schemas.Enum, 0, len(schema.Enums.Enums))

	for _, v := range schema.Enums.Enums {
		enums = append(enums, v)
	}
	sort.Slice(enums, func(i, j int) bool {
		return enums[i].Name < enums[j].Name
	})

	return enums, nil
}
