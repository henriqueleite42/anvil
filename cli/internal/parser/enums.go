package parser

import (
	"fmt"

	"github.com/anvil/anvil/internal/formatter"
	"github.com/anvil/anvil/internal/hashing"
	"github.com/anvil/anvil/schemas"
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
		if self.baseRef != "" {
			ref = self.baseRef + "." + ref
		}
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

		var nameString *string = nil
		nameAny, ok := valuesVMap["Name"]
		if ok {
			localNameString, ok := nameAny.(string)
			if !ok {
				return "", fmt.Errorf("fail to parse \"%s.%s.Values.%s.Name\" to `string`", i.path, i.k, valuesK)
			}
			nameString = &localNameString
		}

		valueString, ok := valuesVMap["Value"].(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Values.%s.Value\" to `string`", i.path, i.k, valuesK)
		}

		values = append(values, &schemas.EnumValue{
			Name:  nameString,
			Value: valueString,
		})
	}

	// TODO make it dynamic based on the metadata
	dbType := formatter.PascalToSnake(i.k) + "_enum"

	rootNode, err := getRootNode(i.path)
	if err != nil {
		return "", err
	}

	enum := &schemas.Enum{
		Ref:          ref,
		OriginalPath: self.getPath(fmt.Sprintf("%s.Enums.%s", i.path, i.k)),
		Name:         i.k,
		DbType:       dbType,
		RootNode:     rootNode,
		Type:         schemas.EnumType(typeString),
		Values:       values,
	}

	stateHash, err := hashing.Struct(enum)
	if err != nil {
		return "", fmt.Errorf("fail to get enum \"%s.%s\" state hash", i.path, i.k)
	}

	enum.StateHash = stateHash
	self.schema.Enums.Enums[refHash] = enum

	return refHash, nil
}

func (self *anvToAnvpParser) enums(file map[string]any) error {
	enumsSchema, ok := file["Enums"]
	if !ok {
		return nil
	}

	fullPath := self.getPath("Enums")

	enumsMap, ok := enumsSchema.(map[string]any)
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
