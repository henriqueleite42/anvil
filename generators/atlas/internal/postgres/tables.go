package postgres

import (
	"fmt"
	"sort"
	"strings"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *hclFile) resolveTables(schema *schemas.Schema) error {
	if schema.Domain == "" {
		return fmt.Errorf("no domain specified")
	}
	if schema.Entities == nil || schema.Entities.Entities == nil {
		return fmt.Errorf("no entities specified")
	}
	if schema.Types == nil || schema.Types.Types == nil {
		return fmt.Errorf("no types specified")
	}

	sortedEntities := []*SortedByOrder{}
	for k, v := range schema.Entities.Entities {
		sortedEntities = append(sortedEntities, &SortedByOrder{
			Order: v.Order,
			Key:   k,
		})
	}
	sort.Slice(sortedEntities, func(i, j int) bool {
		return sortedEntities[i].Order < sortedEntities[j].Order
	})

	tables := []string{}
	for _, sortedEntity := range sortedEntities {
		k := sortedEntity.Key
		v := schema.Entities.Entities[k]

		if v.Name == "" {
			return fmt.Errorf("missing \"Name\" for table \"%s\"", k)
		}

		schemaName, err := self.resolveSchema(v)
		if err != nil {
			return err
		}
		dbSchema := "schema." + schemaName

		sortedColumns := []*SortedByOrder{}
		for k, v := range v.Columns {
			sortedColumns = append(sortedColumns, &SortedByOrder{
				Order: v.Order,
				Key:   k,
			})
		}
		sort.Slice(sortedColumns, func(i, j int) bool {
			return sortedColumns[i].Order < sortedColumns[j].Order
		})

		columnsArr := []string{}
		for _, sortedColumn := range sortedColumns {
			kk := sortedColumn.Key
			vv := v.Columns[kk]

			columnTypeType, ok := schema.Types.Types[vv.TypeHash]
			if !ok {
				return fmt.Errorf("type \"%s\" not found for table \"%s\"", vv.TypeHash, v.Name)
			}

			var columnType string
			if columnTypeType.Type == schemas.TypeType_Enum {
				enumName, err := self.resolveEnum(schema, dbSchema, columnTypeType.EnumHash)
				if err != nil {
					return err
				}
				columnType = "enum." + enumName
			} else {
				columnType = fmt.Sprintf("sql(\"%s\")", *columnTypeType.DbType)
			}

			var autoIncrement string
			if columnTypeType.AutoIncrement {
				autoIncrement = "		auto_increment = true\n"
			}

			var nullable string
			if columnTypeType.Optional {
				autoIncrement = "		null = true\n"
			}

			var defaultVal string
			if columnTypeType.Default != nil {
				autoIncrement = fmt.Sprintf("		default = sql(\"%s\")\n", *columnTypeType.Default)
			}

			column := fmt.Sprintf(`	column "%s" {
		type = %s
%s%s%s	}`, vv.DbName, columnType, autoIncrement, nullable, defaultVal)

			columnsArr = append(columnsArr, column)
		}
		columns := strings.Join(columnsArr, "\n")

		var primaryKeyColumns string
		if v.PrimaryKey != nil && v.PrimaryKey.ColumnsHashes != nil {
			primaryKeyColumnsArr := []string{}
			for _, vv := range v.PrimaryKey.ColumnsHashes {
				column, ok := v.Columns[vv]
				if !ok {
					return fmt.Errorf("column \"%s\" not found for primary key of table \"%s\"", vv, v.Name)
				}
				primaryKeyColumnsArr = append(primaryKeyColumnsArr, "			column."+column.DbName)
			}
			primaryKeyColumns = strings.Join(primaryKeyColumnsArr, "\n")
		}

		var indexes string
		if v.Indexes != nil {
			indexesArr := []string{}
			for _, vv := range v.Indexes {
				columnsArr := []string{}
				for _, vvv := range vv.ColumnsHashes {
					column, ok := v.Columns[vvv]
					if !ok {
						return fmt.Errorf("column \"%s\" not found for index of table \"%s\"", vvv, v.Name)
					}
					columnsArr = append(columnsArr, fmt.Sprintf("			column.%s", column.DbName))
				}
				columns := strings.Join(columnsArr, "\n")

				var unique string
				if vv.Unique {
					unique = "		unique = true\n"
				}

				index := fmt.Sprintf(`	index "%s" {
		columns = [
%s
		]
%s	}`, vv.ConstraintName, columns, unique)

				indexesArr = append(indexesArr, index)
			}

			indexes = strings.Join(indexesArr, "\n") + "\n"
		}

		var foreignKeys string
		if v.ForeignKeys != nil {
			foreignKeysArr := []string{}
			for _, vv := range v.ForeignKeys {
				refTable, ok := schema.Entities.Entities[vv.RefTableHash]
				if !ok {
					return fmt.Errorf("table \"%s\" not found for foreign key of table \"%s\"", vv.RefTableHash, v.Name)
				}

				columnsArr := []string{}
				for _, vvv := range vv.ColumnsHashes {
					column, ok := v.Columns[vvv]
					if !ok {
						return fmt.Errorf("column \"%s\" not found for index of table \"%s\"", vvv, v.Name)
					}
					columnsArr = append(columnsArr, fmt.Sprintf("			column.%s", column.DbName))
				}
				columns := strings.Join(columnsArr, "\n")

				refColumnsArr := []string{}
				for _, vvv := range vv.RefColumnsHashes {
					column, ok := refTable.Columns[vvv]
					if !ok {
						return fmt.Errorf("column \"%s\" not found for ref column of table \"%s\"", vvv, v.Name)
					}
					refColumnsArr = append(refColumnsArr, fmt.Sprintf("			table.%s.column.%s", refTable.DbName, column.DbName))
				}
				refColumns := strings.Join(refColumnsArr, "\n")

				var onUpdate string
				if vv.OnUpdate != nil {
					onUpdate = fmt.Sprintf("		on_update = %s\n", *vv.OnUpdate)
				}

				var onDelete string
				if vv.OnDelete != nil {
					onDelete = fmt.Sprintf("		on_delete = %s\n", *vv.OnDelete)
				}

				index := fmt.Sprintf(`	foreign_key "%s" {
		columns = [
%s
		]
		ref_columns = [
%s
		]
%s%s	}`, vv.ConstraintName, columns, refColumns, onUpdate, onDelete)

				foreignKeysArr = append(foreignKeysArr, index)
			}

			foreignKeys = strings.Join(foreignKeysArr, "\n") + "\n"
		}

		table := fmt.Sprintf(`table "%s" {
	schema = %s
%s
	primary_key {
		columns = [
%s
		]
	}
%s%s}`, v.DbName, dbSchema, columns, primaryKeyColumns, indexes, foreignKeys)

		tables = append(tables, table)
	}

	self.tables = strings.Join(tables, "\n")

	return nil
}
