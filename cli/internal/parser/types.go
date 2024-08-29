package parser

import (
	"fmt"

	"github.com/anvil/anvil/internal/hashing"
	"github.com/anvil/anvil/internal/schema"
)

type AllowedRefs struct {
	All          bool
	Relationship bool
	Type         bool
	Enum         bool
}

func (self *anvToAnvpParser) resolveType(allowedRefs AllowedRefs, path string, k string, v any) (string, error) {
	if self.schema.Types == nil {
		self.schema.Types = &schema.Types{}
	}
	if self.schema.Types.Types == nil {
		self.schema.Types.Types = map[string]*schema.Type{}
	}

	originalPath := path + "." + k
	originalPathHash := hashing.String(originalPath)

	_, ok := self.schema.Types.Types[originalPathHash]
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

	typeTypeAny, ok := vMap["Type"]
	if !ok {
		return "", fmt.Errorf("\"Type\" is a required property to \"%s.%s\"", path, k)
	}
	typeTypeString, ok := typeTypeAny.(string)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s.Type\" to `string`", path, k)
	}
	typeType, ok := schema.ToTypeType(typeTypeString)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s.Type\" to `TypeType`", path, k)
	}

	var confidentiality schema.TypeConfidentiality = schema.TypeConfidentiality_Low
	confidentialityAny, ok := vMap["Confidentiality"]
	if ok {
		confidentialityString, ok := confidentialityAny.(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Confidentiality\" to `string`", path, k)
		}
		confidentiality, ok = schema.ToTypeConfidentiality(confidentialityString)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Confidentiality\" to `TypeConfidentiality`", path, k)
		}
	}

	var optional bool
	optionalAny, ok := vMap["Optional"]
	if ok {
		optionalBool, ok := optionalAny.(bool)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Optional\" to `bool`", path, k)
		}
		optional = optionalBool
	}

	var format *string = nil
	formatAny, ok := vMap["Format"]
	if ok {
		formatString, ok := formatAny.(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Format\" to `string`", path, k)
		}
		format = &formatString
	}

	var validate []string = nil
	validateAny, ok := vMap["Validate"]
	if ok {
		validateArrAny, ok := validateAny.([]any)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Validate\" to `[]any`", path, k)
		}
		validateArr := make([]string, len(validateArrAny)-1)
		for kk, vv := range validateArrAny {
			vString, ok := vv.(string)
			if !ok {
				return "", fmt.Errorf("fail to parse \"%s.%s.Validate.%d\" to `string`", path, k, kk)
			}
			validateArr = append(validateArr, vString)
		}
		validate = validateArr
	}

	var dbType *string = nil
	dbTypeAny, ok := vMap["DbType"]
	if ok {
		dbTypeString, ok := dbTypeAny.(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.DbType\" to `string`", path, k)
		}
		dbType = &dbTypeString
	}

	var childTypesHashes []string = nil
	propertiesAny, ok := vMap["Properties"]
	if ok {
		if typeType != schema.TypeType_ListMap &&
			typeType != schema.TypeType_Map &&
			typeType != schema.TypeType_MapStringMap {
			return "", fmt.Errorf("Type \"%s.%s\" cannot have property \"Properties\". Only types with map \"Type\" can.", path, k)
		}

		propertiesMap, ok := propertiesAny.(map[string]any)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Properties\" to `map[string]any`", path, k)
		}

		typesHashes := []string{}

		for kk, vv := range propertiesMap {
			// Not including "Properties" is intentional to make it smaller and only contain relevant data
			typeHash, err := self.resolveType(allowedRefs, path+"."+k, kk, vv)
			if err != nil {
				return "", err
			}
			typesHashes = append(typesHashes, typeHash)
		}

		childTypesHashes = typesHashes
	} else if typeType == schema.TypeType_ListMap ||
		typeType == schema.TypeType_Map ||
		typeType == schema.TypeType_MapStringMap {
		return "", fmt.Errorf("Type \"%s.%s\" must have property \"Properties\". All types with map \"Type\" must.", path, k)
	}

	var enumHash *string = nil
	valuesAny, ok := vMap["Values"]
	if ok {
		if typeType != schema.TypeType_Enum &&
			typeType != schema.TypeType_ListEnum {
			return "", fmt.Errorf("Type \"%s.%s\" cannot have property \"Values\". Only types with enum \"Type\" can.", path, k)
		}

		valuesHash, err := self.resolveEnum(path+"."+k, "Values", valuesAny)
		if err != nil {
			return "", err
		}

		enumHash = &valuesHash
	} else if typeType == schema.TypeType_Enum ||
		typeType == schema.TypeType_ListEnum {
		return "", fmt.Errorf("Type \"%s.%s\" must have property \"Values\". All types with enum \"Type\" must.", path, k)
	}

	rootNode, err := getRootNode(path)
	if err != nil {
		return "", err
	}

	schemaTypes := &schema.Type{
		Name:             k,
		RootNode:         rootNode,
		OriginalPath:     originalPath,
		Confidentiality:  confidentiality,
		Optional:         optional,
		Format:           format,
		Validate:         validate,
		Type:             typeType,
		DbType:           dbType,
		ChildTypesHashes: childTypesHashes,
		EnumHash:         enumHash,
	}

	stateHash, err := hashing.Struct(schemaTypes)
	if err != nil {
		return "", fmt.Errorf("fail to get import \"%s\" state hash", originalPath)
	}

	schemaTypes.StateHash = stateHash
	self.schema.Types.Types[originalPathHash] = schemaTypes

	return originalPathHash, nil
}

func (self *anvToAnvpParser) types(file map[string]any) error {
	typesSchema, ok := file["Types"]
	if !ok {
		return nil
	}

	fullPath := self.getPath("Types")

	typesMap, ok := typesSchema.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `map[string]any`", fullPath)
	}

	for k, v := range typesMap {
		_, err := self.resolveType(AllowedRefs{
			Relationship: true,
		}, fullPath, k, v)
		if err != nil {
			return err
		}
	}

	return nil
}
