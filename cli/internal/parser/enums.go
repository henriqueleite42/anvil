package parser

import (
	"fmt"
	"sort"

	"github.com/henriqueleite42/anvil/language-helpers/golang/hashing"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *anvToAnvpParser) enums(curDomain string, file map[string]any) error {
	enumsAny, ok := file["Enums"]
	if !ok {
		return nil
	}

	path := curDomain + ".Enums"

	enumsMap, ok := enumsAny.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `map[string]any`", path)
	}

	if self.schema.Enums == nil {
		self.schema.Enums = &schemas.Enums{
			Enums: map[string]*schemas.Enum{},
		}
	}

	for k, v := range enumsMap {
		ref := self.getRef(curDomain, "Enums."+k)
		refHash := hashing.String(ref)

		_, ok := self.schema.Enums.Enums[refHash]
		if ok {
			return fmt.Errorf("fail to parse duplicated enum \"%s.%s\"", path, k)
		}

		vMap, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("fail to parse \"%s.%s\" to `map[string]any`", path, k)
		}

		var description *string
		descriptionAny, ok := vMap["Description"]
		if ok {
			descriptionString, ok := descriptionAny.(string)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.%s.Description\" to `string`", path, k)
			}
			description = &descriptionString
		}

		database := true
		databaseAny, ok := vMap["Database"]
		if ok {
			databaseBoolean, ok := databaseAny.(bool)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.%s.Database\" to `bool`", path, k)
			}
			database = databaseBoolean
		}

		typeAny, ok := vMap["Type"]
		if !ok {
			return fmt.Errorf("\"Type\" is a required property to \"%s.%s\"", path, k)
		}
		typeString, ok := typeAny.(string)
		if !ok {
			return fmt.Errorf("fail to parse \"%s.%s.Type\" to `string`", path, k)
		}

		valuesAny, ok := vMap["Values"]
		if !ok {
			return fmt.Errorf("\"Values\" is a required property to \"%s.%s\"", path, k)
		}
		valuesMap, ok := valuesAny.(map[string]any)
		if !ok {
			return fmt.Errorf("fail to parse \"%s.%s.Values\" to `map[string]any`", path, k)
		}

		values := []*schemas.EnumValue{}
		validateDuplicatedValues := make(map[string]bool, len(valuesMap))
		validateDuplicatedNames := make(map[string]bool, len(valuesMap))
		validateDuplicatedIndexes := make(map[uint]bool, len(valuesMap))
		for valuesK, valuesV := range valuesMap {
			valuesVMap, ok := valuesV.(map[string]any)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.%s.Values.%s\" to `map[string]any`", path, k, valuesK)
			}

			var nameString string
			nameAny, ok := valuesVMap["Name"]
			if ok {
				localNameString, ok := nameAny.(string)
				if !ok {
					return fmt.Errorf("fail to parse \"%s.%s.Values.%s.Name\" to `string`", path, k, valuesK)
				}
				nameString = localNameString
			} else {
				nameString = valuesK
			}

			valueAny, ok := valuesVMap["Value"]
			if !ok {
				return fmt.Errorf("\"%s.%s.Values.%s.Value\" is required", path, k, valuesK)
			}
			valueString, ok := valueAny.(string)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.%s.Values.%s.Value\" to `string`", path, k, valuesK)
			}

			indexAny, ok := valuesVMap["Index"]
			if !ok {
				return fmt.Errorf("\"%s.%s.Values.%s.Index\" is required", path, k, valuesK)
			}
			indexInt, ok := indexAny.(int)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.%s.Values.%s.Index\" to `int`", path, k, valuesK)
			}
			indexUint := uint(indexInt)

			var deprecated bool
			deprecatedAny, ok := valuesVMap["Deprecated"]
			if ok {
				deprecatedBool, ok := deprecatedAny.(bool)
				if !ok {
					return fmt.Errorf("fail to parse \"%s.%s.Values.%s.Deprecated\" to `bool`", path, k, valuesK)
				}
				deprecated = deprecatedBool
			}

			values = append(values, &schemas.EnumValue{
				Name:       nameString,
				Value:      valueString,
				Index:      indexUint,
				Deprecated: deprecated,
			})
		}
		for _, v := range values {
			if validateDuplicatedValues[v.Value] {
				return fmt.Errorf("duplicated enum value \"%s\" for \"%s.%s.Values.%s\"", v.Value, path, k, v.Name)
			}
			if validateDuplicatedNames[v.Name] {
				return fmt.Errorf("duplicated enum name \"%s\" for \"%s.%s.Values.%s\"", v.Name, path, k, v.Name)
			}
			if validateDuplicatedIndexes[v.Index] {
				return fmt.Errorf("duplicated enum index \"%d\" for \"%s.%s.Values.%s\"", v.Index, path, k, v.Name)
			}

			validateDuplicatedValues[v.Value] = true
			validateDuplicatedNames[v.Name] = true
			validateDuplicatedIndexes[v.Index] = true
		}
		sort.Slice(values, func(i, j int) bool {
			return values[i].Index < values[j].Index
		})

		// Adds "Enum" suffix to avoid conflicts
		dbType := self.formatToEntitiesNamingCase(k + "Enum")

		rootNode, err := getRootNode(path)
		if err != nil {
			return err
		}

		enum := &schemas.Enum{
			Ref:          ref,
			OriginalPath: fmt.Sprintf("%s.%s", path, k),
			Domain:       curDomain,
			Name:         k,
			DbName:       dbType,
			DbType:       dbType,
			RootNode:     rootNode,
			Description:  description,
			Database:     database,
			Type:         schemas.EnumType(typeString),
			Values:       values,
		}

		stateHash, err := hashing.Struct(enum)
		if err != nil {
			return fmt.Errorf("fail to get state hash for \"%s.%s\"", path, k)
		}

		enum.StateHash = stateHash
		self.schema.Enums.Enums[refHash] = enum
	}

	return nil
}
