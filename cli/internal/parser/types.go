package parser

import (
	"fmt"
	"sort"

	"github.com/henriqueleite42/anvil/language-helpers/golang/hashing"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

type AllowedRefsWhenInvalid string

func (self *anvToAnvpParser) resolveType(i *resolveInput) (string, error) {
	if self.schema.Types == nil {
		self.schema.Types = &schemas.Types{}
	}
	if self.schema.Types.Types == nil {
		self.schema.Types.Types = map[string]*schemas.Type{}
	}

	ref := self.getDeepRef(i.curDomain, i.ref, i.k)
	refHash := hashing.String(ref)

	_, ok := self.schema.Types.Types[refHash]
	if ok {
		return refHash, nil
	}

	rootNode, err := getRootNode(i.path)
	if err != nil {
		return "", err
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

	var transform []string = nil
	transformAny, ok := vMap["Transform"]
	if ok {
		transformArrAny, ok := transformAny.([]any)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Transform\" to `[]any`", i.path, i.k)
		}
		transformArr := []string{}
		for kk, vv := range transformArrAny {
			vString, ok := vv.(string)
			if !ok {
				return "", fmt.Errorf("fail to parse \"%s.%s.Transform.%d\" to `string`", i.path, i.k, kk)
			}
			transformArr = append(transformArr, vString)
		}
		transform = transformArr
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

	var childTypes []*schemas.TypeChild = nil

	propertiesAny, ok := vMap["Properties"]
	if ok {
		if typeType != schemas.TypeType_Map {
			return "", fmt.Errorf("type \"%s.%s\" cannot have property \"Properties\". Only types with map \"Type\" can", i.path, i.k)
		}

		propertiesMap, ok := propertiesAny.(map[string]any)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Properties\" to `map[string]any`", i.path, i.k)
		}

		childTypes = make([]*schemas.TypeChild, 0, len(propertiesMap))
		for kk, vv := range propertiesMap {
			typeHash, err := self.resolveType(&resolveInput{
				curDomain:  i.curDomain,
				namePrefix: i.namePrefix + i.k,
				path:       fmt.Sprintf("%s.%s.Properties", i.path, i.k),
				ref:        ref,
				k:          kk,
				v:          vv,
			})
			if err != nil {
				return "", err
			}
			childTypes = append(childTypes, &schemas.TypeChild{
				PropName: &kk,
				TypeHash: typeHash,
			})
		}
		sort.Slice(childTypes, func(i, j int) bool {
			return *childTypes[i].PropName < *childTypes[j].PropName
		})
	} else if typeType == schemas.TypeType_Map {
		return "", fmt.Errorf("type \"%s.%s\" must have property \"Properties\". All types with map \"Type\" must", i.path, i.k)
	}

	itemsAny, ok := vMap["Items"]
	if ok {
		if typeType != schemas.TypeType_List {
			return "", fmt.Errorf("type \"%s.%s\" cannot have property \"Items\". Only types with list \"Type\" can", i.path, i.k)
		}

		itemsMap, ok := itemsAny.(map[string]any)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Items\" to `map[string]any`", i.path, i.k)
		}

		kk := i.k + "Item"
		typeHash, err := self.resolveType(&resolveInput{
			curDomain:  i.curDomain,
			namePrefix: i.namePrefix,
			path:       fmt.Sprintf("%s.%s.Items", i.path, i.k),
			ref:        ref,
			k:          kk,
			v:          itemsMap,
		})
		if err != nil {
			return "", err
		}

		if childTypes == nil {
			childTypes = make([]*schemas.TypeChild, 0, 1)
		}

		childTypes = append(childTypes, &schemas.TypeChild{
			TypeHash: typeHash,
		})
	} else if typeType == schemas.TypeType_List {
		return "", fmt.Errorf("type \"%s.%s\" must have property \"Items\". All types with list \"Type\" must", i.path, i.k)
	}

	var enumHash *string = nil
	enumRefAny, ok := vMap["EnumRef"]
	if ok {
		if typeType != schemas.TypeType_Enum {
			return "", fmt.Errorf("type \"%s.%s\" cannot have property \"EnumRef\". Only types with enum \"Type\" can", i.path, i.k)
		}
		enumRefString, ok := enumRefAny.(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.EnumRef\" to `string`", i.path, i.k)
		}
		enumRefString = self.anvRefToAnvpRef(
			i.curDomain,
			enumRefString,
		)

		enumHashVar := hashing.String(enumRefString)
		enumHash = &enumHashVar
	}
	if typeType == schemas.TypeType_Enum && enumHash == nil {
		return "", fmt.Errorf("type \"%s.%s\" must have property \"EnumRef\". All types with enum \"Type\" must", i.path, i.k)
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
		if self.schema.Enums == nil || self.schema.Enums.Enums == nil {
			return "", fmt.Errorf("something went wrong when parsing the enum of \"%s.%s\": no enums parsed", i.path, i.k)
		}
		enum := self.schema.Enums.Enums[*enumHash]

		if enum == nil {
			return "", fmt.Errorf("something went wrong when parsing the enum of \"%s.%s\": enum notfound", i.path, i.k)
		}

		dbType = &enum.DbType
	}

	var dbName *string = nil
	dbNameAny, ok := vMap["DbName"]
	if ok {
		dbNameString, ok := dbNameAny.(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.DbName\" to `string`", i.path, i.k)
		}
		dbName = &dbNameString
	} else if dbType != nil {
		r := self.formatToEntitiesNamingCase(i.k)
		dbName = &r
	}

	name := i.k
	if typeType == schemas.TypeType_Map {
		name = i.namePrefix + i.k
	}

	schemaTypes := &schemas.Type{
		Ref:             ref,
		OriginalPath:    fmt.Sprintf("%s.%s", i.path, i.k),
		Name:            name,
		RootNode:        rootNode,
		Domain:          i.curDomain,
		Confidentiality: confidentiality,
		Optional:        optional,
		Format:          format,
		Validate:        validate,
		Transform:       transform,
		AutoIncrement:   autoIncrement,
		Default:         defaultV,
		Type:            typeType,
		DbName:          dbName,
		DbType:          dbType,
		ChildTypes:      childTypes,
		EnumHash:        enumHash,
	}

	stateHash, err := hashing.Struct(schemaTypes)
	if err != nil {
		return "", fmt.Errorf("fail to get import \"%s.%s\" state hash", i.path, i.k)
	}
	schemaTypes.StateHash = stateHash

	self.schema.Types.Types[refHash] = schemaTypes

	return refHash, nil
}

func (self *anvToAnvpParser) types(curDomain string, file map[string]any) error {
	typesSchema, ok := file["Types"]
	if !ok {
		return nil
	}

	path := curDomain + ".Types"

	typesMap, ok := typesSchema.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `map[string]any`", path)
	}

	for k, v := range typesMap {
		_, err := self.resolveType(&resolveInput{
			curDomain: curDomain,
			path:      path + ".Types",
			ref:       "Types",
			k:         k,
			v:         v,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
