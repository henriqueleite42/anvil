package parse

import (
	"errors"
	"strconv"

	"github.com/anuntech/hephaestus/cmd/types"
)

func getEntityIndexes(key string, data any) ([]*types.Entity_Index, error) {
	valSlice := data.([]any)

	indexes := []*types.Entity_Index{}
	for k, v := range valSlice {
		vMap, ok := v.(map[string]any)
		if !ok {
			return nil, errors.New("fail to parse Entities." + key + ".Indexes." + strconv.Itoa(k))
		}

		columns := []string{}
		if val, ok := vMap["Columns"]; ok {
			valSlice := val.([]any)
			for _, v := range valSlice {
				columns = append(columns, v.(string))
			}
		}

		var unique bool
		if val, ok := vMap["Unique"]; ok {
			valBool, ok := val.(bool)
			if !ok {
				return nil, errors.New("fail to parse Entities." + strconv.Itoa(k) + ".ColumnsCase")
			}
			unique = valBool
		}

		indexes = append(indexes, &types.Entity_Index{
			Columns: columns,
			Unique:  unique,
		})
	}

	return indexes, nil
}

func Entities(s *types.Schema, yaml map[string]any) error {
	entities, ok := yaml["Entities"]
	if !ok {
		return nil
	}

	entitiesAny := entities.(map[string]any)

	var schema *string = nil
	if val, ok := entitiesAny["Schema"]; ok {
		valString := val.(string)
		schema = &valString
	}

	var columnsCase *string = nil
	if val, ok := entitiesAny["ColumnsCase"]; ok {
		valString, ok := val.(types.TextCase)
		if !ok {
			return errors.New("fail to parse Entities.ColumnsCase")
		}
		columnsCase = &valString
	}

	tables := map[string]*types.Entity{}
	tablesMap := entitiesAny["Tables"].(map[string]any)

	for k, v := range tablesMap {
		vMap := v.(map[string]any)

		var name *string = nil
		if val, ok := vMap["Name"]; ok {
			valString := val.(string)
			name = &valString
		}

		var columns map[string]*types.Field
		if val, ok := vMap["Columns"]; ok {
			valMap := val.(map[string]any)
			localColumns, err := resolveField(s, valMap)
			if err != nil {
				return errors.New("fail to parse Entities." + k + ".Columns")
			}
			columns = localColumns
		} else {
			return errors.New("fail to parse Entities." + k + ".Columns")
		}

		var indexes []*types.Entity_Index
		if val, ok := vMap["Indexes"]; ok {
			localIndexes, err := getEntityIndexes(k, val)
			if err != nil {
				return errors.New("fail to parse Entities." + k + ".Indexes")
			}
			indexes = localIndexes
		}

		primaryKeys := []string{}
		if val, ok := vMap["PrimaryKeys"]; ok {
			valSlice := val.([]any)
			for _, v := range valSlice {
				primaryKeys = append(primaryKeys, v.(string))
			}
		}

		tables[k] = &types.Entity{
			Name:        name,
			PrimaryKeys: primaryKeys,
			Columns:     columns,
			Indexes:     indexes,
		}

	}

	s.Entities = &types.Entities{
		Schema:      schema,
		ColumnsCase: columnsCase,
		Tables:      tables,
	}

	return nil
}
