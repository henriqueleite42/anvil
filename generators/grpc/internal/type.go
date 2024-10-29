package internal

import (
	"fmt"
	"strings"

	"github.com/henriqueleite42/anvil/generators/grpc/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/formatter"
	"github.com/henriqueleite42/anvil/language-helpers/golang/grpc"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *parserManager) resolveTypeProp(t *schemas.Type, rootDomain string) (string, error) {
	var typeString string
	if t.Type == schemas.TypeType_String {
		typeString = "string"
	}
	if t.Type == schemas.TypeType_Bytes {
		typeString = "bytes"
	}
	if t.Type == schemas.TypeType_Int ||
		t.Type == schemas.TypeType_Int8 ||
		t.Type == schemas.TypeType_Int16 ||
		t.Type == schemas.TypeType_Int32 {
		typeString = "int32"
	}
	if t.Type == schemas.TypeType_Int64 {
		typeString = "int64"
	}
	if t.Type == schemas.TypeType_Uint ||
		t.Type == schemas.TypeType_Uint8 ||
		t.Type == schemas.TypeType_Uint16 ||
		t.Type == schemas.TypeType_Uint32 {
		typeString = "uint32"
	}
	if t.Type == schemas.TypeType_Uint64 {
		typeString = "uint64"
	}
	if t.Type == schemas.TypeType_Float ||
		t.Type == schemas.TypeType_Float32 {
		typeString = "float"
	}
	if t.Type == schemas.TypeType_Float64 {
		typeString = "double"
	}
	if t.Type == schemas.TypeType_Bool {
		typeString = "bool"
	}
	if t.Type == schemas.TypeType_Timestamp {
		self.grpcTypesParser[t.Domain].imports.AddImport("google/protobuf/timestamp.proto", nil)
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
		if t.ChildTypes == nil {
			return "", fmt.Errorf("type \"%s\" is missing prop \"ChildTypes\"", t.Ref)
		}

		resolvedType, err := self.resolveType(t, rootDomain)
		if err != nil {
			return "", err
		}

		typeString = resolvedType.Name
	}
	if t.Type == schemas.TypeType_List {
		if t.ChildTypes == nil {
			return "", fmt.Errorf("type \"%s\" is missing prop \"ChildTypes\"", t.Name)
		}
		if len(t.ChildTypes) != 1 {
			return "", fmt.Errorf("type \"%s.ChildTypes\" has more than 1 item in the list. It must have exactly one item", t.Name)
		}

		childTypeRef := t.ChildTypes[0].TypeHash
		if childTypeRef == "" {
			return "", fmt.Errorf("type \"%s.ChildTypes\" must have exactly one item", t.Name)
		}

		childType, ok := self.schema.Types.Types[childTypeRef]
		if !ok {
			return "", fmt.Errorf("fail to resolve child type \"%s\" for type \"%s\"", childTypeRef, t.Name)
		}

		typeName, err := self.resolveTypeProp(childType, rootDomain)
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

func (self *parserManager) resolveType(t *schemas.Type, rootDomain string) (*templates.ProtofileTemplInputType, error) {
	if existentType, ok := self.typesToAvoidDuplication[t.Ref]; ok {
		return existentType, nil
	}

	if t.Type != schemas.TypeType_Map {
		return nil, fmt.Errorf("\"%s\" type must be a Map", t.Ref)
	}

	protoTypeName, err := grpc.GetProtoTypeName(t)
	if err != nil {
		return nil, err
	}

	result := &templates.ProtofileTemplInputType{
		Name:  protoTypeName,
		Props: make([]*templates.ProtofileTemplInputTypeProp, 0, len(t.ChildTypes)),
	}

	for k, v := range t.ChildTypes {
		childType, ok := self.schema.Types.Types[v.TypeHash]
		if !ok {
			return nil, fmt.Errorf("child type \"%s\" not found for type \"%s\"", v.TypeHash, t.Name)
		}

		propType, err := self.resolveTypeProp(childType, rootDomain)
		if err != nil {
			return nil, err
		}

		if childType.Domain != rootDomain {
			domainKebab := formatter.PascalToKebab(childType.Domain)
			self.grpcTypesParser[rootDomain].imports.AddImport(domainKebab+".proto", nil)
		}

		result.Props = append(result.Props, &templates.ProtofileTemplInputTypeProp{
			Name: *v.PropName,
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

	if t.RootNode == "Types" || t.RootNode == "Usecase" {
		self.grpcTypesParser[t.Domain].types = append(self.grpcTypesParser[t.Domain].types, result)
	}
	if t.RootNode == "Events" {
		self.grpcTypesParser[t.Domain].events = append(self.grpcTypesParser[t.Domain].events, result)
	}
	if t.RootNode == "Entities" {
		self.grpcTypesParser[t.Domain].entities = append(self.grpcTypesParser[t.Domain].entities, result)
	}

	return result, nil
}
