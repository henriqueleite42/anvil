package internal

import (
	"fmt"
	"strings"

	"github.com/henriqueleite42/anvil/cli/schemas"
)

func (self *protoFile) resolveMsgPropType(schema *schemas.Schema, sourceTypeRef string) (string, error) {
	refType := schema.Types.Types[sourceTypeRef]
	if refType == nil {
		return "", fmt.Errorf("type \"%s\" notfound", sourceTypeRef)
	}

	var typeString string
	if refType.Type == schemas.TypeType_String {
		typeString = "string"
	}
	if refType.Type == schemas.TypeType_Int {
		typeString = "int32"
	}
	if refType.Type == schemas.TypeType_Timestamp {
		self.imports["import \"google/protobuf/timestamp.proto\";"] = true
		typeString = "google.protobuf.Timestamp"
	}
	if refType.Type == schemas.TypeType_Enum {
		if refType.EnumHash == nil {
			return "", fmt.Errorf("type \"%s\" is missing prop \"EnumHash\"", sourceTypeRef)
		}
		if schema.Enums == nil || schema.Enums.Enums == nil {
			return "", fmt.Errorf("missing schema enums, but one of the props requires it")
		}

		enum, ok := schema.Enums.Enums[*refType.EnumHash]
		if !ok {
			return "", fmt.Errorf("enum \"%s\" notfound", *refType.EnumHash)
		}

		typeString = self.resolveEnum(enum)
	}
	if refType.Type == schemas.TypeType_Map {
		if refType.ChildTypesHashes == nil {
			return "", fmt.Errorf("type \"%s\" is missing prop \"ChildTypesHashes\"", sourceTypeRef)
		}

		typeName, err := self.resolveMsg(schema, sourceTypeRef)
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

		typeName, err := self.resolveMsgPropType(schema, refType.ChildTypesHashes[0])
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

func (self *protoFile) resolveMsg(schema *schemas.Schema, sourceTypeRef string) (string, error) {
	refType := schema.Types.Types[sourceTypeRef]
	if refType == nil {
		return "", fmt.Errorf("fail to find type: \"%s\"", sourceTypeRef)
	}

	if refType.Type != schemas.TypeType_Map {
		return "", fmt.Errorf("\"%s\" type must be a Map", refType.Name)
	}

	props := []string{}

	for k, v := range refType.ChildTypesHashes {
		propType, err := self.resolveMsgPropType(schema, v)
		if err != nil {
			return "", err
		}

		props = append(
			props,
			fmt.Sprintf(
				"	%s %s = %d;",
				propType,
				schema.Types.Types[v].Name,
				k,
			),
		)
	}

	self.messages[refType.Name] = fmt.Sprintf(`message %s {
%s
}`, refType.Name, strings.Join(props, "\n"))

	return refType.Name, nil
}
