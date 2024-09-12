package parser

import (
	"fmt"
	"strings"

	"github.com/henriqueleite42/anvil/cli/schemas"
	"github.com/henriqueleite42/anvil/generators/grpc-client-go/internal/templates"
)

func (self *parserManager) toType(t *schemas.Type) (*templates.TemplType, error) {
	if existentType, ok := self.types[t.Name]; ok {
		return existentType, nil
	}

	result := &templates.TemplType{
		OriginalType: t.Type,
	}

	if t.Type == schemas.TypeType_String {
		result.Name = "string"
	}
	if t.Type == schemas.TypeType_Int {
		result.Name = "int32"
	}
	if t.Type == schemas.TypeType_Float {
		result.Name = "float32"
	}
	if t.Type == schemas.TypeType_Bool {
		result.Name = "bool"
	}
	if t.Type == schemas.TypeType_Timestamp {
		self.importsContract["time"] = true
		result.Name = "time.Time"
	}
	if t.Type == schemas.TypeType_Enum {
		if t.EnumHash == nil {
			return nil, fmt.Errorf("enum \"%s\" not found", *t.EnumHash)
		}

		schemaEnum := self.schema.Enums.Enums[*t.EnumHash]
		enum, err := self.toEnum(schemaEnum)
		if err != nil {
			return nil, err
		}

		result.Name = enum.Name
	}
	if t.Type == schemas.TypeType_List {
		if t.ChildTypesHashes == nil {
			return nil, fmt.Errorf("ChildTypesHashes for \"%s\" not found", t.Name)
		}
		if len(t.ChildTypesHashes) != 1 {
			return nil, fmt.Errorf("ChildTypesHashes for \"%s\" must have exactly one item", t.Name)
		}

		childType, ok := self.schema.Types.Types[t.ChildTypesHashes[0]]
		if !ok {
			return nil, fmt.Errorf("type \"%s\" not found", t.ChildTypesHashes[0])
		}

		resolvedChildType, err := self.toType(childType)
		if err != nil {
			return nil, err
		}

		result.Name = "[]" + resolvedChildType.Name
	}
	if t.Type == schemas.TypeType_Map {
		biggest := 0
		types := make([]*schemas.Type, 0, len(t.ChildTypesHashes))
		for _, v := range t.ChildTypesHashes {
			sType, ok := self.schema.Types.Types[v]
			if !ok {
				return nil, fmt.Errorf("type \"%s\" not found", v)
			}

			types = append(types, sType)

			newLen := len(sType.Name)
			if newLen > biggest {
				biggest = newLen
			}
		}

		result.Name = t.Name
		result.Props = make([]*templates.TemplTypeProp, 0, len(types))

		for _, v := range types {
			targetLen := biggest - len(v.Name)

			propType, err := self.toType(v)
			if err != nil {
				return nil, err
			}

			resultPropType := propType.Name
			if propType.OriginalType == schemas.TypeType_Map {
				resultPropType = "*" + resultPropType
			}

			result.Props = append(result.Props, &templates.TemplTypeProp{
				Name:    v.Name,
				Spacing: strings.Repeat(" ", targetLen),
				Type:    resultPropType,
			})
		}

		self.types[t.Name] = result
	}

	if t.Optional && t.Type != schemas.TypeType_Map && t.Type != schemas.TypeType_List {
		result.Name = "*" + result.Name
	}

	return result, nil
}
