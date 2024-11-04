package postgres

import (
	"fmt"
	"sort"

	"github.com/henriqueleite42/anvil/generators/atlas/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func resolveEntities(schema *schemas.AnvpSchema) ([]*templates.HclTemplInputEntity, error) {
	result := make([]*templates.HclTemplInputEntity, 0, len(schema.Entities.Entities))

	for _, v := range schema.Entities.Entities {
		columns := make([]*templates.HclTemplInputEntityColumn, 0, len(v.Columns))
		primaryKeys := make([]*templates.HclTemplInputEntityPrimaryColumn, 0, len(v.PrimaryKey.ColumnsHashes))
		indexes := make([]*templates.HclTemplInputEntityIndex, 0, len(v.Indexes))
		foreignKeys := make([]*templates.HclTemplInputEntityForeignKey, 0, len(v.ForeignKeys))

		for _, vv := range v.Columns {
			cType, ok := schema.Types.Types[vv.TypeHash]
			if !ok {
				return nil, fmt.Errorf("type \"%s\" not found", vv.TypeHash)
			}
			if cType.DbName == nil {
				return nil, fmt.Errorf("type \"%s\" must have \"DbName\"", cType.Name)
			}

			var columnType string
			if cType.Type == schemas.TypeType_Enum {
				if cType.EnumHash == nil {
					return nil, fmt.Errorf("type \"%s\" must have `EnumHash`", cType.Name)
				}
				if schema.Enums == nil || schema.Enums.Enums == nil {
					return nil, fmt.Errorf("enum for type \"%s\" not found: schema has no enums", cType.Name)
				}
				enum, ok := schema.Enums.Enums[*cType.EnumHash]
				if !ok {
					return nil, fmt.Errorf("enum for type \"%s\" not found: enum with hash \"%s\" not found", cType.Name, *cType.EnumHash)
				}

				columnType = fmt.Sprintf("enum.%s", enum.DbName)
			} else {
				if cType.DbType == nil {
					return nil, fmt.Errorf("\"%s\" must have a `DbType`", cType.Name)
				}

				columnType = fmt.Sprintf("sql(\"%s\")", *cType.DbType)
			}

			var cDefault *string
			if cType.Default != nil {
				val := fmt.Sprintf("sql(\"%s\")", *cType.Default)
				cDefault = &val
			}

			columns = append(columns, &templates.HclTemplInputEntityColumn{
				Name:          *cType.DbName,
				Type:          columnType,
				Default:       cDefault,
				Optional:      cType.Optional,
				AutoIncrement: cType.AutoIncrement,
				Order:         vv.Order,
			})
		}
		sort.Slice(columns, func(i, j int) bool {
			return columns[i].Order < columns[j].Order
		})

		for _, vv := range v.PrimaryKey.ColumnsHashes {
			column, ok := v.Columns[vv]
			if !ok {
				return nil, fmt.Errorf("fail to get column with hash \"%s\" to parse primary key for table \"%s\"", vv, v.Name)
			}
			cType, ok := schema.Types.Types[column.TypeHash]
			if !ok {
				return nil, fmt.Errorf("type \"%s\" not found", column.TypeHash)
			}
			if cType.DbName == nil {
				return nil, fmt.Errorf("type \"%s\" must have \"DbName\"", cType.Name)
			}
			primaryKeys = append(primaryKeys, &templates.HclTemplInputEntityPrimaryColumn{
				DbName: "column." + *cType.DbName,
				Order:  column.Order,
			})
		}
		sort.Slice(primaryKeys, func(i, j int) bool {
			return primaryKeys[i].Order < primaryKeys[j].Order
		})

		for _, vv := range v.Indexes {
			columnsNames := make([]string, 0, len(vv.ColumnsHashes))
			for _, vv := range vv.ColumnsHashes {
				column, ok := v.Columns[vv]
				if !ok {
					return nil, fmt.Errorf("fail to get column with hash \"%s\" for table \"%s\"", vv, v.Name)
				}
				cType, ok := schema.Types.Types[column.TypeHash]
				if !ok {
					return nil, fmt.Errorf("type \"%s\" not found", column.TypeHash)
				}
				if cType.DbName == nil {
					return nil, fmt.Errorf("type \"%s\" must have \"DbName\"", cType.Name)
				}
				columnsNames = append(columnsNames, "column."+*cType.DbName)
			}

			indexes = append(indexes, &templates.HclTemplInputEntityIndex{
				Name:    vv.ConstraintName,
				Columns: columnsNames,
				Unique:  vv.Unique,
			})
		}
		// TODO: Add sort when order added to Indexes

		for _, vv := range v.ForeignKeys {
			refTableEntity, ok := schema.Entities.Entities[vv.RefTableHash]
			if !ok {
				return nil, fmt.Errorf("fail to get ref table with hash \"%s\"", vv.RefTableHash)
			}

			columns := make([]string, 0, len(vv.ColumnsHashes))
			for _, vv := range vv.ColumnsHashes {
				column, ok := v.Columns[vv]
				if !ok {
					return nil, fmt.Errorf("fail to get column with hash \"%s\" for table \"%s\"", vv, v.Name)
				}
				cType, ok := schema.Types.Types[column.TypeHash]
				if !ok {
					return nil, fmt.Errorf("type \"%s\" not found", column.TypeHash)
				}
				if cType.DbName == nil {
					return nil, fmt.Errorf("type \"%s\" must have \"DbName\"", cType.Name)
				}
				columns = append(columns, "column."+*cType.DbName)
			}
			refColumns := make([]string, 0, len(vv.RefColumnsHashes))
			for _, vv := range vv.RefColumnsHashes {
				column, ok := refTableEntity.Columns[vv]
				if !ok {
					return nil, fmt.Errorf("fail to get table with hash \"%s\"", vv)
				}
				cType, ok := schema.Types.Types[column.TypeHash]
				if !ok {
					return nil, fmt.Errorf("type \"%s\" not found", column.TypeHash)
				}
				if cType.DbName == nil {
					return nil, fmt.Errorf("type \"%s\" must have \"DbName\"", cType.Name)
				}
				refColumns = append(refColumns, "table."+refTableEntity.DbName+".column."+*cType.DbName)
			}

			foreignKeys = append(foreignKeys, &templates.HclTemplInputEntityForeignKey{
				RefTable:   refTableEntity.DbName,
				Name:       vv.ConstraintName,
				Columns:    columns,
				RefColumns: refColumns,
				OnUpdate:   vv.OnUpdate,
				OnDelete:   vv.OnDelete,
			})
		}
		// TODO: Add sort when order added to ForeignKeys

		result = append(result, &templates.HclTemplInputEntity{
			DbName:      v.DbName,
			Columns:     columns,
			PrimaryKeys: primaryKeys,
			Indexes:     indexes,
			ForeignKeys: foreignKeys,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].DbName < result[j].DbName
	})

	return result, nil
}
