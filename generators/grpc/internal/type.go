package internal

import (
	"fmt"
	"strings"

	"github.com/henriqueleite42/anvil/generators/grpc/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *parser) resolveTypeProp(t *schemas.Type) (string, error) {
	var typeString string
	if t.Type == schemas.TypeType_String {
		typeString = "string"
	}
	if t.Type == schemas.TypeType_Int {
		typeString = "int32"
	}
	if t.Type == schemas.TypeType_Float {
		typeString = "float"
	}
	if t.Type == schemas.TypeType_Bool {
		typeString = "bool"
	}
	if t.Type == schemas.TypeType_Timestamp {
		self.imports["google/protobuf/timestamp.proto"] = true
		typeString = "google.protobuf.Timestamp"
	}
	if t.Type == schemas.TypeType_Enum {
		if self.schema.Enums == nil || self.schema.Enums.Enums == nil {
			return "", fmt.Errorf("missing schema enums, but one of the props requires it")
		}
		if t.EnumHash == nil {
			return "", fmt.Errorf("type \"%s\" is missing prop \"EnumHash\"", t.Ref)
		}

		enum, ok := self.schema.Enums.Enums[*t.EnumHash]
		if !ok {
			return "", fmt.Errorf("enum \"%s\" notfound", *t.EnumHash)
		}

		enumResolved := self.resolveEnum(enum)

		typeString = enumResolved.Name
	}
	if t.Type == schemas.TypeType_Map {
		if t.ChildTypesHashes == nil {
			return "", fmt.Errorf("type \"%s\" is missing prop \"ChildTypesHashes\"", t.Ref)
		}

		resolvedType, err := self.resolveType(t)
		if err != nil {
			return "", err
		}

		typeString = resolvedType.Name
	}
	if t.Type == schemas.TypeType_List {
		if t.ChildTypesHashes == nil {
			return "", fmt.Errorf("type \"%s\" is missing prop \"ChildTypesHashes\"", t.Name)
		}
		if len(t.ChildTypesHashes) != 1 {
			return "", fmt.Errorf("type \"%s.ChildTypesHashes\" has more than 1 item in the list. It must have exactly one item.", t.Name)
		}

		childTypeRef := t.ChildTypesHashes[0]
		if childTypeRef == "" {
			return "", fmt.Errorf("type \"%s.ChildTypesHashes\" must have exactly one item.", t.Name)
		}

		childType, ok := self.schema.Types.Types[childTypeRef]
		if !ok {
			return "", fmt.Errorf("fail to resolve child type \"%s\" for type \"%s\".", childTypeRef, t.Name)
		}

		typeName, err := self.resolveTypeProp(childType)
		if err != nil {
			return "", err
		}

		typeString = fmt.Sprintf("repeated %s", typeName)
	}
	if typeString == "" {
		return "", fmt.Errorf("fail to find type for \"%s\"", t.Type)
	}

	if t.Optional && t.Type != schemas.TypeType_List {
		typeString = "optional " + typeString
	}

	return typeString, nil
}

func (self *parser) resolveType(t *schemas.Type) (*templates.ProtofileTemplInputType, error) {
	if existentType, ok := self.typesToAvoidDuplication[t.Ref]; ok {
		return existentType, nil
	}

	if t.Type != schemas.TypeType_Map {
		return nil, fmt.Errorf("\"%s\" type must be a Map", t.Name)
	}

	result := &templates.ProtofileTemplInputType{
		Name:  t.Name,
		Props: make([]*templates.ProtofileTemplInputTypeProp, 0, len(t.ChildTypesHashes)),
	}

	for k, v := range t.ChildTypesHashes {
		childType, ok := self.schema.Types.Types[v]
		if !ok {
			return nil, fmt.Errorf("child type \"%s\" not found for type \"%s\"", v, t.Name)
		}

		propType, err := self.resolveTypeProp(childType)
		if err != nil {
			return nil, err
		}

		result.Props = append(result.Props, &templates.ProtofileTemplInputTypeProp{
			Name: childType.Name,
			Type: propType,
			Idx:  k + 1,
		})
	}

	biggestName := 0
	biggestType := 0
	for _, v := range result.Props {
		newLenName := len(v.Name)
		if newLenName > biggestName {
			biggestName = newLenName
		}

		newLenType := len(v.Type)
		if newLenType > biggestType {
			biggestType = newLenType
		}
	}

	for _, v := range result.Props {
		v.Spacing1 = strings.Repeat(" ", biggestType-len(v.Type))
		v.Spacing2 = strings.Repeat(" ", biggestName-len(v.Name))
	}

	self.typesToAvoidDuplication[t.Ref] = result
	self.types = append(self.types, result)

	return result, nil
}