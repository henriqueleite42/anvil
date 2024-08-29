package parser

import (
	"fmt"

	"github.com/anvil/anvil/internal/hashing"
	"github.com/anvil/anvil/internal/schema"
)

func (self *anvToAnvpParser) resolveEnum(path string, k string, v any) (string, error) {
	if self.schema.Enums == nil {
		self.schema.Enums = &schema.Enums{}
	}
	if self.schema.Enums.Enums == nil {
		self.schema.Enums.Enums = map[string]*schema.Enum{}
	}

	originalPath := path + "." + k
	originalPathHash := hashing.String(originalPath)

	_, ok := self.schema.Enums.Enums[originalPathHash]
	if ok {
		return originalPathHash, nil
	}

	vMap, ok := v.(map[string]any)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s\" to `map[string]any`", path, k)
	}

	// TODO
	_, ok = vMap["$ref"]
	if ok {
		return "", nil
	}

	typeAny, ok := vMap["Type"]
	if !ok {
		return "", fmt.Errorf("\"Type\" is a required property to \"%s.%s\"", path, k)
	}
	typeString, ok := typeAny.(string)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s.Type\" to `string`", path, k)
	}

	valuesAny, ok := vMap["Values"]
	if !ok {
		return "", fmt.Errorf("\"Values\" is a required property to \"%s.%s\"", path, k)
	}
	valuesMap, ok := valuesAny.(map[string]any)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s.Values\" to `map[string]any`", path, k)
	}

	values := []*schema.EnumValue{}
	for valuesK, valuesV := range valuesMap {
		valuesVMap, ok := valuesV.(map[string]any)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Values.%s\" to `map[string]any`", path, k, valuesK)
		}

		var nameString *string = nil
		nameAny, ok := valuesVMap["Name"]
		if ok {
			localNameString, ok := nameAny.(string)
			if !ok {
				return "", fmt.Errorf("fail to parse \"%s.%s.Values.%s.Name\" to `string`", path, k, valuesK)
			}
			nameString = &localNameString
		}

		valueString, ok := valuesVMap["Value"].(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Values.%s.Value\" to `string`", path, k, valuesK)
		}

		values = append(values, &schema.EnumValue{
			Name:  nameString,
			Value: valueString,
		})
	}

	rootNode, err := getRootNode(path)
	if err != nil {
		return "", err
	}

	enum := &schema.Enum{
		Name:         k,
		RootNode:     rootNode,
		OriginalPath: originalPath,
		Type:         schema.EnumType(typeString),
		Values:       values,
	}

	stateHash, err := hashing.Struct(enum)
	if err != nil {
		return "", fmt.Errorf("fail to get enum \"%s\" state hash", originalPath)
	}

	enum.StateHash = stateHash
	self.schema.Enums.Enums[originalPathHash] = enum

	return originalPathHash, nil
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
		_, err := self.resolveEnum(fullPath, k, v)
		if err != nil {
			return err
		}
	}

	return nil
}
