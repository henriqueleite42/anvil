package contract

import (
	"fmt"
	"strings"

	"github.com/henriqueleite42/anvil/cli/schemas"
	"github.com/henriqueleite42/anvil/generators/grpc-client-go/internal/spacer"
)

func (self *contractFile) resolveType(sourceTypeRef string, parentTypeName string) (string, error) {
	refType := self.schema.Types.Types[sourceTypeRef]
	if refType == nil {
		return "", fmt.Errorf("type \"%s\" notfound", sourceTypeRef)
	}

	var typeString string
	if refType.Type == schemas.TypeType_String {
		typeString = "string"
	}
	if refType.Type == schemas.TypeType_Int {
		typeString = "int"
	}
	if refType.Type == schemas.TypeType_Float {
		typeString = "float32"
	}
	if refType.Type == schemas.TypeType_Bool {
		typeString = "bool"
	}
	if refType.Type == schemas.TypeType_Timestamp {
		self.imports["	\"time\""] = true
		typeString = "time.Time"
	}
	if refType.Type == schemas.TypeType_Enum {
		if refType.EnumHash == nil {
			return "", fmt.Errorf("type \"%s\" is missing prop \"EnumHash\"", sourceTypeRef)
		}
		if self.schema.Enums == nil || self.schema.Enums.Enums == nil {
			return "", fmt.Errorf("missing schema enums, but one of the props requires it")
		}

		enum, ok := self.schema.Enums.Enums[*refType.EnumHash]
		if !ok {
			return "", fmt.Errorf("enum \"%s\" notfound", *refType.EnumHash)
		}

		enumName, err := self.resolveEnum(enum)
		if err != nil {
			return "", err
		}

		typeString = enumName
	}
	if refType.Type == schemas.TypeType_Map {
		if refType.ChildTypesHashes == nil {
			return "", fmt.Errorf("type \"%s\" is missing prop \"ChildTypesHashes\"", sourceTypeRef)
		}

		typeName, err := self.resolveTypeMap(sourceTypeRef, "")
		if err != nil {
			return "", err
		}

		typeString = typeName
	}
	if refType.Type == schemas.TypeType_List {
		if refType.ChildTypesHashes == nil {
			return "", fmt.Errorf("type \"%s\" is missing prop \"ChildTypesHashes\"", sourceTypeRef)
		}
		if len(refType.ChildTypesHashes) != 1 {
			return "", fmt.Errorf("type \"%s.ChildTypesHashes\" has more than 1 item in the list. It must have exactly one item.", sourceTypeRef)
		}

		typeName, err := self.resolveType(refType.ChildTypesHashes[0], parentTypeName)
		if err != nil {
			return "", err
		}

		typeString = fmt.Sprintf("[]%s", typeName)
	}
	if typeString == "" {
		return "", fmt.Errorf("fail to find type for \"%s\"", refType.Type)
	}

	return typeString, nil
}

func (self *contractFile) resolveTypeMap(sourceTypeRef string, parentTypeName string) (string, error) {
	refType := self.schema.Types.Types[sourceTypeRef]
	if refType == nil {
		return "", fmt.Errorf("fail to find type: \"%s\"", sourceTypeRef)
	}

	if _, ok := self.types[parentTypeName+refType.Name]; ok {
		return parentTypeName + refType.Name, nil
	}

	if refType.Type != schemas.TypeType_Map {
		return "", fmt.Errorf("\"%s\" type must be a Map", refType.Name)
	}

	props, err := spacer.Space(
		refType.ChildTypesHashes,
		func(v string) ([]string, error) {
			propType, err := self.resolveType(v, refType.Name)
			if err != nil {
				return nil, err
			}

			typeType := self.schema.Types.Types[v]

			return []string{
				parentTypeName + typeType.Name,
				propType,
			}, nil
		},
		func(s []string, i int) (string, error) {
			targetLen := i - len(s[0])
			return fmt.Sprintf(
				"	%s %s",
				s[0]+strings.Repeat(" ", targetLen),
				s[1],
			), nil
		},
	)
	if err != nil {
		return "", err
	}

	self.types[parentTypeName+refType.Name] = &ItemWithOrder{
		Order: len(self.types),
		Value: fmt.Sprintf(`type %s struct {
%s
}`, refType.Name, strings.Join(props, "\n")),
	}

	return "*" + refType.Name, nil
}
