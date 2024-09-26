package internal

import (
	"fmt"
	"strings"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *protoFile) resolveMsgPropType(parentType *schemas.Type, sourceTypeRef string) (string, error) {
	refType := self.schema.Types.Types[sourceTypeRef]
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
	if refType.Type == schemas.TypeType_Float {
		typeString = "float"
	}
	if refType.Type == schemas.TypeType_Bool {
		typeString = "bool"
	}
	if refType.Type == schemas.TypeType_Timestamp {
		self.imports["import \"google/protobuf/timestamp.proto\";"] = true
		typeString = "google.protobuf.Timestamp"
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

		typeString = self.resolveEnum(enum)
	}
	if refType.Type == schemas.TypeType_Map {
		if refType.ChildTypesHashes == nil {
			return "", fmt.Errorf("type \"%s\" is missing prop \"ChildTypesHashes\"", sourceTypeRef)
		}

		typeName, err := self.resolveMsg(sourceTypeRef)
		if err != nil {
			return "", err
		}

		typeString = typeName
	}
	if refType.Type == schemas.TypeType_List {
		if parentType != nil && parentType.Type == schemas.TypeType_List {
			return "", fmt.Errorf("fail to parse \"%s\": grpc doesn't support lists of lists", refType.Name)
		}

		if refType.ChildTypesHashes == nil {
			return "", fmt.Errorf("type \"%s\" is missing prop \"ChildTypesHashes\"", refType.Name)
		}
		if len(refType.ChildTypesHashes) != 1 {
			return "", fmt.Errorf("type \"%s.ChildTypesHashes\" has more than 1 item in the list. It must have exactly one item.", refType.Name)
		}

		typeName, err := self.resolveMsgPropType(refType, refType.ChildTypesHashes[0])
		if err != nil {
			return "", err
		}

		typeString = fmt.Sprintf("repeated %s", typeName)
	}
	if typeString == "" {
		return "", fmt.Errorf("fail to find type for \"%s\"", refType.Type)
	}

	return typeString, nil
}

func (self *protoFile) resolveMsg(sourceTypeRef string) (string, error) {
	refType := self.schema.Types.Types[sourceTypeRef]
	if refType == nil {
		return "", fmt.Errorf("fail to find type: \"%s\"", sourceTypeRef)
	}

	if _, ok := self.messages[refType.Name]; ok {
		return refType.Name, nil
	}

	if refType.Type != schemas.TypeType_Map {
		return "", fmt.Errorf("\"%s\" type must be a Map", refType.Name)
	}

	props := []string{}

	for k, v := range refType.ChildTypesHashes {
		propType, err := self.resolveMsgPropType(nil, v)
		if err != nil {
			return "", err
		}

		typeType := self.schema.Types.Types[v]

		var optional string
		if typeType.Optional && typeType.Type != schemas.TypeType_List {
			optional = "optional "
		}

		props = append(
			props,
			fmt.Sprintf(
				"	%s%s %s = %d;",
				optional,
				propType,
				typeType.Name,
				k+1,
			),
		)
	}

	self.messages[refType.Name] = fmt.Sprintf(`message %s {
%s
}`, refType.Name, strings.Join(props, "\n"))

	return refType.Name, nil
}
