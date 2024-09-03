package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/henriqueleite42/anvil/cli/internal/formatter"
	"github.com/henriqueleite42/anvil/cli/internal/hashing"
	"github.com/henriqueleite42/anvil/cli/schemas"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func (self *anvToAnvpParser) resolveEntity(i *resolveInput) (string, error) {
	if self.schema.Entities == nil {
		self.schema.Entities = &schemas.Entities{}
	}
	if self.schema.Entities.Entities == nil {
		self.schema.Entities.Entities = map[string]*schemas.Entity{}
	}

	ref := "Entities." + i.k
	if i.ref != "" {
		ref = i.ref + "." + ref
	}
	refHash := hashing.String(ref)

	path := fmt.Sprintf("%s.%s", i.path, i.k)

	_, ok := self.schema.Entities.Entities[refHash]
	if ok {
		return refHash, nil
	}

	vMap, ok := i.v.(map[string]any)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s\" to `map[string]any`", i.path, i.k)
	}

	refAny, ok := vMap["$ref"]
	if ok {
		refString, ok := refAny.(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.$ref\" to `string`", i.path, i.k)
		}
		return hashing.String(refString), nil
	}

	var tableSchema *string = nil
	tableSchemaAny, ok := vMap["Schema"]
	if ok {
		tableSchemaString, ok := tableSchemaAny.(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Schema\" to `string`", i.path, i.k)
		}
		tableSchema = &tableSchemaString
	}

	var tableName string
	tableNameAny, ok := vMap["Name"]
	if ok {
		tableNameString, ok := tableNameAny.(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Name\" to `string`", i.path, i.k)
		}
		tableName = tableNameString
	}
	if tableName == "" {
		tableName = i.k
	}

	rootNode, err := getRootNode(i.path)
	if err != nil {
		return "", err
	}

	// Also create a type for the entity, not only for the columns
	entityType := &schemas.Type{
		Ref:              ref,
		OriginalPath:     path,
		Name:             i.k,
		RootNode:         rootNode,
		ChildTypesHashes: []string{},
		Type:             schemas.TypeType_Map,
		// TODO Implement metadata to control this
		// or tell by the columns: If one is medium, them set to medium, etc
		Confidentiality: schemas.TypeConfidentiality_Low,
	}

	columnsAny, ok := vMap["Columns"]
	if !ok {
		return "", fmt.Errorf("\"Columns\" is a required property to \"%s.%s\"", i.path, i.k)
	}
	columnsArr, ok := columnsAny.(map[string]any)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s.Columns\" to `map[string]any`", i.path, i.k)
	}
	columns := map[string]*schemas.EntityColumn{}
	for kk, vv := range columnsArr {
		vvMap, ok := vv.(map[string]any)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Columns.%s\" to `map[string]any`", i.path, i.k, kk)
		}

		var columnName string
		columnNameAny, ok := vvMap["Name"]
		if ok {
			columnNameString, ok := columnNameAny.(string)
			if !ok {
				return "", fmt.Errorf("fail to parse \"%s.%s.Columns.%s.Name\" to `string`", i.path, i.k, kk)
			}
			columnName = columnNameString
		} else if self.schema.Entities != nil &&
			self.schema.Entities.Metadata != nil &&
			self.schema.Entities.Metadata.ColumnsCase != nil {
			if *self.schema.Entities.Metadata.ColumnsCase == schemas.ColumnsCase_Snake {
				columnName = formatter.PascalToSnake(kk)
			} else if *self.schema.Entities.Metadata.ColumnsCase == schemas.ColumnsCase_Camel {
				columnName = formatter.PascalToCamel(kk)
			} else {
				columnName = kk
			}
		} else {
			columnName = kk
		}

		columnPath := fmt.Sprintf("%s.%s.Columns.%s", i.path, i.k, kk)
		columnRef := ref + "." + kk

		var typeHash string
		typeHash, err := self.resolveType(&resolveInput{
			path: fmt.Sprintf("%s.%s.Columns", i.path, i.k),
			ref:  ref,
			k:    kk,
			v:    vv,
		})
		if err != nil {
			return "", err
		}

		columnRefHash := hashing.String(columnRef)
		column := &schemas.EntityColumn{
			Ref:          columnRef,
			OriginalPath: columnPath,
			Name:         kk,
			ColumnName:   columnName,
			TypeHash:     typeHash,
		}

		stateHash, err := hashing.Struct(column)
		if err != nil {
			return "", fmt.Errorf("fail to get state hash for %s.%s.Columns.%s", i.path, i.k, kk)
		}
		column.StateHash = stateHash

		columns[columnRefHash] = column
		entityType.ChildTypesHashes = append(entityType.ChildTypesHashes, columnRefHash)
	}

	entityTypeStateHash, err := hashing.Struct(entityType)
	if err != nil {
		return "", fmt.Errorf("fail to get state hash for \"%s.%s\" (type)", i.path, i.k)
	}
	entityType.StateHash = entityTypeStateHash

	self.schema.Types.Types[refHash] = entityType

	primaryKeyAny, ok := vMap["PrimaryKey"]
	if !ok {
		return "", fmt.Errorf("\"PrimaryKey\" is a required property to \"%s.%s\"", i.path, i.k)
	}
	primaryKeyArr, ok := primaryKeyAny.([]any)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s.PrimaryKey\" to `[]any`", i.path, i.k)
	}
	pkColumnsHashes := []string{}
	for cni, columnNameAny := range primaryKeyArr {
		columnNameString, ok := columnNameAny.(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.PrimaryKey.%d\" to `[]any`", i.path, i.k, cni)
		}
		columnRef := fmt.Sprintf("Entities.%s.%s", i.k, columnNameString)
		hash := hashing.String(columnRef)
		pkColumnsHashes = append(pkColumnsHashes, hash)
	}
	primaryKey := &schemas.EntityPrimaryKey{
		ConstraintName: tableName + "_pk",
		ColumnsHashes:  pkColumnsHashes, // TODO
	}
	primaryKeyHash, err := hashing.Struct(primaryKey)
	if err != nil {
		return "", fmt.Errorf("fail to get \"%s.%s.PrimaryKey\" state hash", i.path, i.k)
	}
	primaryKey.StateHash = primaryKeyHash

	var indexes map[string]*schemas.EntityIndex = nil
	indexesAny, ok := vMap["Indexes"]
	if ok {
		indexesArr, ok := indexesAny.([]any)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Indexes\" to `[]any`", i.path, i.k)
		}

		indexes = map[string]*schemas.EntityIndex{}

		for kk, vv := range indexesArr {
			vvMap, ok := vv.(map[string]any)
			if !ok {
				return "", fmt.Errorf("fail to parse \"%s.%s.Indexes.%d\" to `map[string]any`", i.path, i.k, kk)
			}

			columnsHashes := []string{}
			columnsAny, ok := vvMap["Columns"]
			if !ok {
				return "", fmt.Errorf("\"Columns\" is a required property to \"%s.%s.Indexes.%d\"", i.path, i.k, kk)
			}
			columnsArr, ok := columnsAny.([]any)
			if !ok {
				return "", fmt.Errorf("fail to parse \"%s.%s.Indexes.%d.Columns\" to `[]any`", i.path, i.k, kk)
			}
			for kkk, vvv := range columnsArr {
				vvvString, ok := vvv.(string)
				if !ok {
					return "", fmt.Errorf("fail to parse \"%s.%s.Indexes.%d.Columns.%d\" to `string`", i.path, i.k, kk, kkk)
				}
				columnRef := fmt.Sprintf("Entities.%s.%s", i.k, vvvString)
				if i.ref != "" {
					columnRef = i.ref + "." + columnRef
				}
				hash := hashing.String(columnRef)
				columnsHashes = append(columnsHashes, hash)
			}

			var name string
			nameAny, ok := vvMap["Name"]
			if ok {
				nameBool, ok := nameAny.(string)
				if !ok {
					return "", fmt.Errorf("fail to parse \"%s.%s.Indexes.%d.Name\" to `string`", i.path, i.k, kk)
				}
				name = nameBool
			} else {
				columnsToCreateName := []string{}
				for _, vvv := range columnsHashes {
					column, ok := columns[vvv]
					if !ok {
						return "", fmt.Errorf("fail to find one of the columns to created database name for \"%s.%s.Indexes.%d\"", i.path, i.k, kk)
					}
					columnsToCreateName = append(columnsToCreateName, column.ColumnName)
				}
				// TODO make it dynamic to match pattern specified in Entities.ColumnsCase (maybe create a Entities.ConstraintCase?)
				name = strings.Join(columnsToCreateName, "_") + "_idx"
			}

			var unique bool
			uniqueAny, ok := vvMap["Unique"]
			if ok {
				uniqueBool, ok := uniqueAny.(bool)
				if !ok {
					return "", fmt.Errorf("fail to parse \"%s.%s.Indexes.%d.Unique\" to `bool`", i.path, i.k, kk)
				}
				unique = uniqueBool
			}

			indexPath := fmt.Sprintf("%s.Indexes.%d", path, kk)
			indexHash := hashing.String(indexPath)

			index := &schemas.EntityIndex{
				OriginalPath:   indexPath,
				ConstraintName: name,
				ColumnsHashes:  columnsHashes,
				Unique:         unique,
			}

			stateHash, err := hashing.Struct(index)
			if err != nil {
				return "", fmt.Errorf("fail to get \"%s.%s.Indexes.%d\" hash", i.path, i.k, kk)
			}
			index.StateHash = stateHash

			indexes[indexHash] = index
		}
	}

	var foreignKeys map[string]*schemas.EntityForeignKey = nil
	foreignKeysAny, ok := vMap["ForeignKeys"]
	if ok {
		foreignKeysArr, ok := foreignKeysAny.([]any)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.ForeignKeys\" to `[]any`", i.path, i.k)
		}

		foreignKeys = map[string]*schemas.EntityForeignKey{}

		for kk, vv := range foreignKeysArr {
			vvMap, ok := vv.(map[string]any)
			if !ok {
				return "", fmt.Errorf("fail to parse \"%s.%s.ForeignKeys.%d\" to `map[string]any`", i.path, i.k, kk)
			}

			columnsNamesForFkName := []string{}
			columnsHashes := []string{}
			columnsAny, ok := vvMap["Columns"]
			if !ok {
				return "", fmt.Errorf("\"Columns\" is a required property to \"%s.%s.ForeignKeys.%d\"", i.path, i.k, kk)
			}
			columnsArr, ok := columnsAny.([]any)
			if !ok {
				return "", fmt.Errorf("fail to parse \"%s.%s.ForeignKeys.%d.Columns\" to `[]any`", i.path, i.k, kk)
			}
			for kkk, vvv := range columnsArr {
				vvvString, ok := vvv.(string)
				if !ok {
					return "", fmt.Errorf("fail to parse \"%s.%s.ForeignKeys.%d.Columns.%d\" to `string`", i.path, i.k, kk, kkk)
				}
				columnRef := fmt.Sprintf("Entities.%s.%s", i.k, vvvString)
				hash := hashing.String(columnRef)

				column := columns[hash]

				if column == nil {
					return "", fmt.Errorf("fail to find column \"%s\" for \"%s.%s.ForeignKeys.%d.Columns.%d\"", vvvString, i.path, i.k, kk, kkk)
				}

				columnsNamesForFkName = append(columnsNamesForFkName, column.ColumnName)
				columnsHashes = append(columnsHashes, hash)
			}

			var refTableHash string

			refColumnsHashes := []string{}
			refColumnsAny, ok := vvMap["RefColumns"]
			if !ok {
				return "", fmt.Errorf("\"RefColumns\" is a required property to \"%s.%s.ForeignKeys.%d\"", i.path, i.k, kk)
			}
			refColumnsArr, ok := refColumnsAny.([]any)
			if !ok {
				return "", fmt.Errorf("fail to parse \"%s.%s.ForeignKeys.%d.RefColumns\" to `[]any`", i.path, i.k, kk)
			}
			if len(refColumnsArr) == 0 {
				return "", fmt.Errorf("at least 1 column is required for \"%s.%s.ForeignKeys.%d.RefColumns\"", i.path, i.k, kk)
			}
			var firstTableRefName string
			for kkk, vvv := range refColumnsArr {
				vvvString, ok := vvv.(string)
				if !ok {
					return "", fmt.Errorf("fail to parse \"%s.%s.ForeignKeys.%d.RefColumns.%d\" to `string`", i.path, i.k, kk, kkk)
				}

				parts := strings.Split(vvvString, ".")
				if len(parts) != 2 {
					return "", fmt.Errorf("You can only reference columns using the pattern \"TableName.ColumnName\" on \"%s.%s.ForeignKeys.%d.RefColumns.%d\" to `string`", i.path, i.k, kk, kkk)
				}

				refTableName := parts[0]
				refColumnName := parts[1]

				if firstTableRefName == "" {
					firstTableRefName = refTableName
				}

				if firstTableRefName != refTableName {
					return "", fmt.Errorf("You can only reference one ref table per foreign key. Error found on \"%s.%s.ForeignKeys.%d.RefColumns.%d\" to `string`", i.path, i.k, kk, kkk)
				}

				if refTableHash == "" {
					refTableHash = hashing.String("Entities." + refTableName)
				}

				refColumn := "Entities." + refTableName + "." + refColumnName

				refColumnHash := hashing.String(refColumn)

				refColumnsHashes = append(refColumnsHashes, refColumnHash)
			}

			var name string
			nameAny, ok := vvMap["Name"]
			if ok {
				nameString, ok := nameAny.(string)
				if !ok {
					return "", fmt.Errorf("fail to parse \"%s.%s.ForeignKeys.%d.Name\" to `string`", i.path, i.k, kk)
				}
				name = nameString
			} else {
				name = fmt.Sprintf("%s_%s_fk", tableName, strings.Join(columnsNamesForFkName, "_"))
			}

			var onDelete *string = nil
			onDeleteAny, ok := vvMap["OnDelete"]
			if ok {
				onDeleteString, ok := onDeleteAny.(string)
				if !ok {
					return "", fmt.Errorf("fail to parse \"%s.%s.ForeignKeys.%d.OnDelete\" to `string`", i.path, i.k, kk)
				}
				onDelete = &onDeleteString
			}

			var onUpdate *string = nil
			onUpdateAny, ok := vvMap["OnUpdate"]
			if ok {
				onUpdateString, ok := onUpdateAny.(string)
				if !ok {
					return "", fmt.Errorf("fail to parse \"%s.%s.ForeignKeys.%d.OnUpdate\" to `string`", i.path, i.k, kk)
				}
				onUpdate = &onUpdateString
			}

			fkPath := fmt.Sprintf("%s.ForeignKeys.%d", i.path, kk)
			fkHash := hashing.String(fkPath)

			fk := &schemas.EntityForeignKey{
				OriginalPath:     fkPath,
				ConstraintName:   name,
				ColumnsHashes:    columnsHashes,
				RefTableHash:     refTableHash,
				RefColumnsHashes: refColumnsHashes,
				OnDelete:         onDelete,
				OnUpdate:         onUpdate,
			}

			stateHash, err := hashing.Struct(fk)
			if err != nil {
				return "", fmt.Errorf("fail to get \"%s.%s.ForeignKeys.%d\" hash", i.path, i.k, kk)
			}
			fk.StateHash = stateHash

			foreignKeys[fkHash] = fk
		}
	}

	entity := &schemas.Entity{
		Ref:          ref,
		OriginalPath: path,
		Name:         i.k,
		RootNode:     rootNode,
		TypeHash:     refHash,
		Schema:       tableSchema,
		TableName:    tableName,
		Columns:      columns,
		PrimaryKey:   primaryKey,
		Indexes:      indexes,
		ForeignKeys:  foreignKeys,
	}

	stateHash, err := hashing.Struct(entity)
	if err != nil {
		return "", fmt.Errorf("fail to get state hash for \"%s.%s\"", i.path, i.k)
	}
	entity.StateHash = stateHash

	self.schema.Entities.Entities[refHash] = entity

	return refHash, nil
}

func (self *anvToAnvpParser) resolveEntitiesMetadata(file map[string]any) error {
	entitiesSchema, ok := file["Entities"]
	if !ok {
		return nil
	}

	path := "Entities"

	entitiesMap, ok := entitiesSchema.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `map[string]any`", path)
	}

	if self.schema.Entities == nil {
		self.schema.Entities = &schemas.Entities{}
	}

	var columnsCase *schemas.ColumnsCase = nil
	columnsCaseAny, ok := entitiesMap["ColumnsCase"]
	if ok {
		columnsCaseString, ok := columnsCaseAny.(string)
		if !ok {
			return fmt.Errorf("fail to parse \"%s.ColumnsCase\" to `string`", path)
		}
		if !ok {
			return fmt.Errorf("fail to parse \"%s.ColumnsCase\" to `string`", path)
		}
		columnsCaseStr, ok := schemas.ToColumnsCase(columnsCaseString)
		if !ok {
			return fmt.Errorf("fail to parse \"%s.ColumnsCase\" to `TypeConfidentiality`", path)
		}
		columnsCase = &columnsCaseStr
	}

	self.schema.Entities.Metadata = &schemas.EntitiesMetadata{
		ColumnsCase: columnsCase,
	}

	return nil
}

func (self *anvToAnvpParser) entities(file map[string]any) error {
	entitiesSchema, ok := file["Entities"]
	if !ok {
		return nil
	}

	path := "Entities"

	entitiesMap, ok := entitiesSchema.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `map[string]any`", path)
	}

	if self.schema.Entities == nil {
		self.schema.Entities = &schemas.Entities{}
	}

	entitiesAny, ok := entitiesMap["Entities"]
	if !ok {
		return fmt.Errorf("\"Entities\" is a required property to \"%s.Entities.Entities\"", path)
	}
	entitiesMap, ok = entitiesAny.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s.Entities.Entities\" to `map[string]any`", path)
	}

	for k, v := range entitiesMap {
		_, err := self.resolveEntity(&resolveInput{
			path: path + ".Entities",
			ref:  "",
			k:    k,
			v:    v,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
