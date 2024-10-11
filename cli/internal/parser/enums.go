package parser

import (
	"fmt"
	"sort"

	"github.com/henriqueleite42/anvil/language-helpers/golang/hashing"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *anvToAnvpParser) resolveEnum(i *resolveInput) (string, error) {
	if self.schema.Enums == nil {
		self.schema.Enums = &schemas.Enums{}
	}
	if self.schema.Enums.Enums == nil {
		self.schema.Enums.Enums = map[string]*schemas.Enum{}
	}

	// Enums ref works a little different,
	// because they also can be created Entities, Usecase, Repository, etc
	// so we use their reference instead of the Enum
	var ref string
	if i.ref != "" {
		ref = i.ref + "." + i.k
	} else {
		ref = "Enums." + i.k
	}
	refHash := hashing.String(ref)

	_, ok := self.schema.Enums.Enums[refHash]
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

	typeAny, ok := vMap["Type"]
	if !ok {
		return "", fmt.Errorf("\"Type\" is a required property to \"%s.%s\"", i.path, i.k)
	}
	typeString, ok := typeAny.(string)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s.Type\" to `string`", i.path, i.k)
	}

	valuesAny, ok := vMap["Values"]
	if !ok {
		return "", fmt.Errorf("\"Values\" is a required property to \"%s.%s\"", i.path, i.k)
	}
	valuesMap, ok := valuesAny.(map[string]any)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s.Values\" to `map[string]any`", i.path, i.k)
	}

	values := []*schemas.EnumValue{}
	for valuesK, valuesV := range valuesMap {
		valuesVMap, ok := valuesV.(map[string]any)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Values.%s\" to `map[string]any`", i.path, i.k, valuesK)
		}

		var nameString string
		nameAny, ok := valuesVMap["Name"]
		if ok {
			localNameString, ok := nameAny.(string)
			if !ok {
				return "", fmt.Errorf("fail to parse \"%s.%s.Values.%s.Name\" to `string`", i.path, i.k, valuesK)
			}
			nameString = localNameString
		} else {
			nameString = valuesK
		}

		valueAny, ok := valuesVMap["Value"]
		if !ok {
			return "", fmt.Errorf("\"%s.%s.Values.%s.Value\" is required", i.path, i.k, valuesK)
		}
		valueString, ok := valueAny.(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Values.%s.Value\" to `string`", i.path, i.k, valuesK)
		}

		indexAny, ok := valuesVMap["Index"]
		if !ok {
			return "", fmt.Errorf("\"%s.%s.Values.%s.Index\" is required", i.path, i.k, valuesK)
		}
		indexInt, ok := indexAny.(int)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Values.%s.Index\" to `int`", i.path, i.k, valuesK)
		}
		indexInt32 := int32(indexInt)

		var deprecated bool
		deprecatedAny, ok := valuesVMap["Deprecated"]
		if ok {
			deprecatedBool, ok := deprecatedAny.(bool)
			if !ok {
				return "", fmt.Errorf("fail to parse \"%s.%s.Values.%s.Deprecated\" to `int32`", i.path, i.k, valuesK)
			}
			deprecated = deprecatedBool
		}

		values = append(values, &schemas.EnumValue{
			Name:       nameString,
			Value:      valueString,
			Index:      indexInt32,
			Deprecated: deprecated,
		})
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i].Index < values[j].Index
	})

	dbType := self.formatToEntitiesNamingCase(i.k + "Enum")

	rootNode, err := getRootNode(i.path)
	if err != nil {
		return "", err
	}

	enum := &schemas.Enum{
		Ref:          ref,
		OriginalPath: fmt.Sprintf("%s.Enums.%s", i.path, i.k),
		Name:         i.k,
		DbName:       dbType,
		DbType:       dbType,
		RootNode:     rootNode,
		Type:         schemas.EnumType(typeString),
		Values:       values,
	}

	stateHash, err := hashing.Struct(enum)
	if err != nil {
		return "", fmt.Errorf("fail to get state hash for \"%s.%s\"", i.path, i.k)
	}

	enum.StateHash = stateHash
	self.schema.Enums.Enums[refHash] = enum

	return refHash, nil
}

func (self *anvToAnvpParser) enums(file map[string]any) error {
	enumsAny, ok := file["Enums"]
	if !ok {
		return nil
	}

	fullPath := "Enums"

	enumsMap, ok := enumsAny.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `map[string]any`", fullPath)
	}

	for k, v := range enumsMap {
		_, err := self.resolveEnum(&resolveInput{
			path: fullPath,
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
