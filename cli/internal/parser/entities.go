package parser

import (
	"fmt"
	"sort"
	"strings"

	"github.com/henriqueleite42/anvil/language-helpers/golang/hashing"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *anvToAnvpParser) resolveEntity(i *resolveInput) (string, error) {
	if self.schema.Entities == nil {
		self.schema.Entities = &schemas.Entities{}
	}
	if self.schema.Entities.Entities == nil {
		self.schema.Entities.Entities = map[string]*schemas.Entity{}
	}

	ref := self.getDeepRef(i.curDomain, i.ref, "Entities."+i.k)
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
		refString = self.anvRefToAnvpRef(
			i.curDomain,
			refString,
		)

		return hashing.String(refString), nil
	}

	var dbSchema *string = nil
	dbSchemaAny, ok := vMap["Schema"]
	if ok {
		dbSchemaString, ok := dbSchemaAny.(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Schema\" to `string`", i.path, i.k)
		}
		dbSchema = &dbSchemaString
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
		// Pluralize and format
		tableName = self.formatToEntitiesNamingCase(i.k + "s")
	}

	rootNode, err := getRootNode(i.path)
	if err != nil {
		return "", err
	}

	// Also create a type for the entity, not only for the columns
	entityType := &schemas.Type{
		Ref:          ref,
		OriginalPath: path,
		Name:         i.k,
		RootNode:     rootNode,
		ChildTypes:   []*schemas.TypeChild{},
		Type:         schemas.TypeType_Map,
		// Entities does not have Confidentiality levels, only their fields
		Confidentiality: schemas.TypeConfidentiality_Low,
	}

	columnsAny, ok := vMap["Columns"]
	if !ok {
		return "", fmt.Errorf("\"Columns\" is a required property to \"%s.%s\"", i.path, i.k)
	}
	columnsMap, ok := columnsAny.(map[string]any)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s.Columns\" to `map[string]any`", i.path, i.k)
	}
	columns := map[string]*schemas.EntityColumn{}

	// Necessary to keep some kind of order
	keys := []string{}
	for key := range columnsMap {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for columnOrder, kk := range keys {
		vv := columnsMap[kk]

		columnPath := fmt.Sprintf("%s.%s.Columns.%s", i.path, i.k, kk)
		columnRef := fmt.Sprintf("%s.%s", ref, kk)

		typeHash, err := self.resolveType(&resolveInput{
			curDomain: i.curDomain,
			path:      fmt.Sprintf("%s.%s.Columns", i.path, i.k),
			ref:       ref,
			k:         kk,
			v:         vv,
		})
		if err != nil {
			return "", err
		}

		columnRefHash := hashing.String(columnRef)
		column := &schemas.EntityColumn{
			Ref:          columnRef,
			OriginalPath: columnPath,
			Order:        uint(columnOrder),
			Name:         kk,
			TypeHash:     typeHash,
		}

		stateHash, err := hashing.Struct(column)
		if err != nil {
			return "", fmt.Errorf("fail to get state hash for %s.%s.Columns.%s", i.path, i.k, kk)
		}
		column.StateHash = stateHash

		columns[columnRefHash] = column
		entityType.ChildTypes = append(entityType.ChildTypes, &schemas.TypeChild{
			PropName: &kk,
			TypeHash: columnRefHash,
		})

		columnOrder++
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
				columnRef := fmt.Sprintf("%s.%s", ref, vvvString)
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
						return "", fmt.Errorf("fail to find one of the columns to create database name for \"%s.%s.Indexes.%d\"", i.path, i.k, kk)
					}
					cType, ok := self.schema.Types.Types[column.TypeHash]
					if !ok {
						return "", fmt.Errorf("fail to find columns type to create database name for \"%s.%s.Indexes.%d.%s\"", i.path, i.k, kk, column.Name)
					}
					if cType.DbName == nil {
						return "", fmt.Errorf("fail to find get DbName to create database name for \"%s.%s.Indexes.%d\"", i.path, i.k, kk)
					}
					columnsToCreateName = append(columnsToCreateName, *cType.DbName)
				}
				name = fmt.Sprintf("%s_%s_idx", tableName, strings.Join(columnsToCreateName, "_"))
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
				Order:          uint(kk),
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
				columnRef := fmt.Sprintf("%s.%s", ref, vvvString)
				hash := hashing.String(columnRef)

				column, ok := columns[hash]
				if !ok {
					return "", fmt.Errorf("fail to find column \"%s\" for \"%s.%s.ForeignKeys.%d.Columns.%d\"", vvvString, i.path, i.k, kk, kkk)
				}

				cType, ok := self.schema.Types.Types[column.TypeHash]
				if !ok {
					return "", fmt.Errorf("fail to find one of the columns types to created database name for \"%s.%s.Indexes.%d\"", i.path, i.k, kk)
				}
				if cType.DbName == nil {
					return "", fmt.Errorf("fail to find get DbName to created database name for \"%s.%s.Indexes.%d\"", i.path, i.k, kk)
				}

				columnsNamesForFkName = append(columnsNamesForFkName, *cType.DbName)
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
					return "", fmt.Errorf("you can only reference columns using the pattern \"TableName.ColumnName\" on \"%s.%s.ForeignKeys.%d.RefColumns.%d\" to `string`", i.path, i.k, kk, kkk)
				}

				refTableName := parts[0]
				refColumnName := parts[1]

				if firstTableRefName == "" {
					firstTableRefName = refTableName
				}

				if firstTableRefName != refTableName {
					return "", fmt.Errorf("you can only reference one ref table per foreign key. Error found on \"%s.%s.ForeignKeys.%d.RefColumns.%d\" to `string`", i.path, i.k, kk, kkk)
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
				Order:            uint(kk),
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

	order := len(self.schema.Entities.Entities)
	entity := &schemas.Entity{
		Ref:          ref,
		OriginalPath: path,
		Name:         i.k,
		RootNode:     rootNode,
		TypeHash:     refHash,
		Order:        uint(order),
		DbSchema:     dbSchema,
		DbName:       tableName,
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

func (self *anvToAnvpParser) resolveEntitiesMetadata(curDomain string, file map[string]any) error {
	entitiesSchema, ok := file["Entities"]
	if !ok {
		return nil
	}

	path := curDomain + ".Entities"

	entitiesMap, ok := entitiesSchema.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `map[string]any`", path)
	}

	if self.schema.Entities == nil {
		self.schema.Entities = &schemas.Entities{}
	}

	var columnsCase *schemas.NamingCase = nil
	columnsCaseAny, ok := entitiesMap["NamingCase"]
	if ok {
		columnsCaseString, ok := columnsCaseAny.(string)
		if !ok {
			return fmt.Errorf("fail to parse \"%s.NamingCase\" to `string`", path)
		}
		if !ok {
			return fmt.Errorf("fail to parse \"%s.NamingCase\" to `string`", path)
		}
		columnsCaseStr, ok := schemas.ToNamingCase(columnsCaseString)
		if !ok {
			return fmt.Errorf("fail to parse \"%s.NamingCase\" to `TypeConfidentiality`", path)
		}
		columnsCase = &columnsCaseStr
	}

	self.schema.Entities.Metadata = &schemas.EntitiesMetadata{
		NamingCase: columnsCase,
	}

	return nil
}

func (self *anvToAnvpParser) entities(curDomain string, file map[string]any) error {
	entitiesSchema, ok := file["Entities"]
	if !ok {
		return nil
	}

	path := curDomain + ".Entities"

	entitiesMap, ok := entitiesSchema.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `map[string]any`", path)
	}

	if self.schema.Entities == nil {
		self.schema.Entities = &schemas.Entities{}
	}
	if self.schema.Entities.Entities == nil {
		self.schema.Entities.Entities = map[string]*schemas.Entity{}
	}

	entitiesAny, ok := entitiesMap["Entities"]
	if !ok {
		return fmt.Errorf("\"Entities\" is a required property to \"%s.Entities.Entities\"", path)
	}
	entitiesMap, ok = entitiesAny.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s.Entities.Entities\" to `map[string]any`", path)
	}

	// Necessary to keep some kind of order
	keys := []string{}
	for key := range entitiesMap {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, k := range keys {
		v := entitiesMap[k]
		_, err := self.resolveEntity(&resolveInput{
			curDomain: curDomain,
			path:      path + ".Entities",
			ref:       "",
			k:         k,
			v:         v,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
