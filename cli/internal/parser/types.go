package parser

import (
	"fmt"

	"github.com/henriqueleite42/anvil/cli/hashing"
	"github.com/henriqueleite42/anvil/cli/schemas"
)

type AllowedRefsWhenInvalid string

func (self *anvToAnvpParser) resolveType(i *resolveInput) (string, error) {
	if self.schema.Types == nil {
		self.schema.Types = &schemas.Types{}
	}
	if self.schema.Types.Types == nil {
		self.schema.Types.Types = map[string]*schemas.Type{}
	}

	// Types ref works a little different,
	// because they also can be created Entities, Usecase, Repository, etc
	// so we use their reference instead of the Type
	var ref string
	if i.ref != "" {
		ref = i.ref + "." + i.k
	} else {
		ref = "Types." + i.k
	}
	refHash := hashing.String(ref)

	_, ok := self.schema.Types.Types[refHash]
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

	typeTypeAny, ok := vMap["Type"]
	if !ok {
		return "", fmt.Errorf("\"Type\" is a required property to \"%s.%s\"", i.path, i.k)
	}
	typeTypeString, ok := typeTypeAny.(string)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s.Type\" to `string`", i.path, i.k)
	}
	typeType, ok := schemas.ToTypeType(typeTypeString)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s.Type\" to `TypeType`", i.path, i.k)
	}

	var confidentiality schemas.TypeConfidentiality = schemas.TypeConfidentiality_Low
	confidentialityAny, ok := vMap["Confidentiality"]
	if ok {
		confidentialityString, ok := confidentialityAny.(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Confidentiality\" to `string`", i.path, i.k)
		}
		confidentiality, ok = schemas.ToTypeConfidentiality(confidentialityString)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Confidentiality\" to `TypeConfidentiality`", i.path, i.k)
		}
	}

	var optional bool
	optionalAny, ok := vMap["Optional"]
	if ok {
		optionalBool, ok := optionalAny.(bool)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Optional\" to `bool`", i.path, i.k)
		}
		optional = optionalBool
	}

	var format *string = nil
	formatAny, ok := vMap["Format"]
	if ok {
		formatString, ok := formatAny.(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Format\" to `string`", i.path, i.k)
		}
		format = &formatString
	}

	var validate []string = nil
	validateAny, ok := vMap["Validate"]
	if ok {
		validateArrAny, ok := validateAny.([]any)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Validate\" to `[]any`", i.path, i.k)
		}
		validateArr := []string{}
		for kk, vv := range validateArrAny {
			vString, ok := vv.(string)
			if !ok {
				return "", fmt.Errorf("fail to parse \"%s.%s.Validate.%d\" to `string`", i.path, i.k, kk)
			}
			validateArr = append(validateArr, vString)
		}
		validate = validateArr
	}

	var autoIncrement bool
	autoIncrementAny, ok := vMap["AutoIncrement"]
	if ok {
		autoIncrementBool, ok := autoIncrementAny.(bool)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Format\" to `bool`", i.path, i.k)
		}
		autoIncrement = autoIncrementBool
	}

	var defaultV *string = nil
	defaultVAny, ok := vMap["Default"]
	if ok {
		defaultVString, ok := defaultVAny.(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Default\" to `string`", i.path, i.k)
		}
		defaultV = &defaultVString
	}

	var childTypesHashes []string = nil
	propertiesAny, ok := vMap["Properties"]
	if ok {
		if typeType != schemas.TypeType_Map {
			return "", fmt.Errorf("Type \"%s.%s\" cannot have property \"Properties\". Only types with map \"Type\" can.", i.path, i.k)
		}

		propertiesMap, ok := propertiesAny.(map[string]any)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Properties\" to `map[string]any`", i.path, i.k)
		}

		typesHashes := []string{}

		for kk, vv := range propertiesMap {
			typeHash, err := self.resolveType(&resolveInput{
				path: fmt.Sprintf("%s.%s.Properties", i.path, i.k),
				ref:  ref,
				k:    kk,
				v:    vv,
			})
			if err != nil {
				return "", err
			}
			typesHashes = append(typesHashes, typeHash)
		}

		childTypesHashes = typesHashes
	} else if typeType == schemas.TypeType_Map {
		return "", fmt.Errorf("Type \"%s.%s\" must have property \"Properties\". All types with map \"Type\" must.", i.path, i.k)
	}

	var enumHash *string = nil
	valuesAny, ok := vMap["Values"]
	if ok {
		if typeType != schemas.TypeType_Enum {
			return "", fmt.Errorf("Type \"%s.%s\" cannot have property \"Values\". Only types with enum \"Type\" can.", i.path, i.k)
		}

		valuesHash, err := self.resolveEnum(&resolveInput{
			path: i.path + ".Values",
			ref:  ref,
			k:    i.k,
			v:    valuesAny,
		})
		if err != nil {
			return "", err
		}

		enumHash = &valuesHash
	}
	if typeType == schemas.TypeType_Enum && enumHash == nil {
		return "", fmt.Errorf("Type \"%s.%s\" must have property \"Values\". All types with enum \"Type\" must.", i.path, i.k)
	}

	var dbType *string = nil
	dbTypeAny, ok := vMap["DbType"]
	if ok {
		dbTypeString, ok := dbTypeAny.(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.DbType\" to `string`", i.path, i.k)
		}
		dbType = &dbTypeString
	} else if typeType == schemas.TypeType_Enum {
		// TODO make it dynamic to match pattern specified in Entities.NamingCase (maybe create a Entities.ConstraintCase?)
		if self.schema.Enums == nil || self.schema.Enums.Enums == nil {
			return "", fmt.Errorf("something went wrong when parsing the enum of \"%s.%s\": no enums parsed.", i.path, i.k)
		}
		enum := self.schema.Enums.Enums[*enumHash]

		if enum == nil {
			return "", fmt.Errorf("something went wrong when parsing the enum of \"%s.%s\": enum notfound", i.path, i.k)
		}

		dbType = &enum.DbType
	}

	rootNode, err := getRootNode(i.path)
	if err != nil {
		return "", err
	}

	schemaTypes := &schemas.Type{
		Ref:              ref,
		OriginalPath:     fmt.Sprintf("%s.%s", i.path, i.k),
		Name:             i.k,
		RootNode:         rootNode,
		Confidentiality:  confidentiality,
		Optional:         optional,
		Format:           format,
		Validate:         validate,
		AutoIncrement:    autoIncrement,
		Default:          defaultV,
		Type:             typeType,
		DbType:           dbType,
		ChildTypesHashes: childTypesHashes,
		EnumHash:         enumHash,
	}

	stateHash, err := hashing.Struct(schemaTypes)
	if err != nil {
		return "", fmt.Errorf("fail to get import \"%s.%s\" state hash", i.path, i.k)
	}
	schemaTypes.StateHash = stateHash

	self.schema.Types.Types[refHash] = schemaTypes

	return refHash, nil
}

func (self *anvToAnvpParser) types(file map[string]any) error {
	typesSchema, ok := file["Types"]
	if !ok {
		return nil
	}

	path := "Types"

	typesMap, ok := typesSchema.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `map[string]any`", path)
	}

	for k, v := range typesMap {
		_, err := self.resolveType(&resolveInput{
			path: "Types.Types",
			ref:  "Types",
			k:    k,
			v:    v,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
