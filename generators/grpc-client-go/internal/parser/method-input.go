package parser

import (
	"fmt"
	"strings"

	"github.com/henriqueleite42/anvil/cli/formatter"
	"github.com/henriqueleite42/anvil/generators/grpc-client-go/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *parserManager) toMethodInput(methodName string, t *schemas.Type) ([]*templates.TemplMethodProp, []string, error) {
	if t.Type != schemas.TypeType_Map {
		return nil, nil, fmt.Errorf("inputs for grpc must be Map Type")
	}

	_, err := self.toGoType(t)
	if err != nil {
		return nil, nil, err
	}

	biggest := 0
	amountOfOptional := 0
	types := []*schemas.Type{}
	for _, v := range t.ChildTypesHashes {
		propType, ok := self.schema.Types.Types[v]
		if !ok {
			return nil, nil, fmt.Errorf("type \"%s\" not found", v)
		}

		types = append(types, propType)

		if len(propType.Name) > biggest {
			biggest = len(propType.Name)
		}

		if propType.Optional {
			amountOfOptional++
		}
	}

	props := []*templates.TemplMethodProp{}
	propsPrepare := []string{}
	for _, propType := range types {
		var value string

		if propType.Type == schemas.TypeType_String ||
			propType.Type == schemas.TypeType_Int ||
			propType.Type == schemas.TypeType_Float ||
			propType.Type == schemas.TypeType_Bool {
			value = "i." + propType.Name
		}
		if propType.Type == schemas.TypeType_Timestamp {
			self.importsImplementation["google.golang.org/protobuf/types/known/timestamppb"] = true

			if propType.Optional {
				varName := formatter.PascalToCamel(propType.Name)
				prepareList, err := self.templateManager.Parse("input-prop-optional", &templates.InputPropOptionalTemplInput{
					VarName:              varName,
					OriginalVariableName: fmt.Sprintf("i.%s", propType.Name),
					Type:                 "*timestamppb.Timestamp",
					ValueToAssign:        fmt.Sprintf("timestamppb.New(*i.%s)", propType.Name),
				})
				if err != nil {
					return nil, nil, err
				}
				propsPrepare = append(propsPrepare, prepareList)
				value = varName
			} else {
				value = fmt.Sprintf("timestamppb.New(i.%s)", propType.Name)
			}
		}
		if propType.Type == schemas.TypeType_Enum {
			if propType.EnumHash == nil {
				return nil, nil, fmt.Errorf("enum \"%s\" not found", *propType.EnumHash)
			}

			schemaEnum := self.schema.Enums.Enums[*propType.EnumHash]
			enum, err := self.toEnum(schemaEnum)
			if err != nil {
				return nil, nil, err
			}

			if propType.Optional {
				varName := formatter.PascalToCamel(propType.Name)
				prepareList, err := self.templateManager.Parse("input-prop-optional", &templates.InputPropOptionalTemplInput{
					VarName:              varName,
					OriginalVariableName: fmt.Sprintf("i.%s", propType.Name),
					Type:                 "*pb." + enum.Name,
					ValueToAssign:        fmt.Sprintf("convert%sToPb(*i.%s)", enum.Name, propType.Name),
					NeedsPointer:         true,
				})
				if err != nil {
					return nil, nil, err
				}
				propsPrepare = append(propsPrepare, prepareList)
				value = varName
			} else {
				value = fmt.Sprintf("convert%sToPb(i.%s)", enum.Name, propType.Name)
			}
		}
		if propType.Type == schemas.TypeType_List {
			if propType.Optional {
				return nil, nil, fmt.Errorf("unable to parse \"%s\": grpc-client-go currently doesn't support optional lists", propType.Name)
			}
			if propType.ChildTypesHashes == nil {
				return nil, nil, fmt.Errorf("ChildTypesHashes for \"%s\" not found", t.Name)
			}
			if len(propType.ChildTypesHashes) != 1 {
				return nil, nil, fmt.Errorf("ChildTypesHashes for \"%s\" must have exactly one item", t.Name)
			}

			childType, ok := self.schema.Types.Types[propType.ChildTypesHashes[0]]
			if !ok {
				return nil, nil, fmt.Errorf("type \"%s\" not found", propType.ChildTypesHashes[0])
			}

			varName := formatter.PascalToCamel(propType.Name)

			if childType.Type == schemas.TypeType_String ||
				childType.Type == schemas.TypeType_Int ||
				childType.Type == schemas.TypeType_Float ||
				childType.Type == schemas.TypeType_Bool {
				value = "i." + propType.Name
			} else {
				var childTypeType string
				var childTypeToAppend string
				if childType.Type == schemas.TypeType_Timestamp {
					self.importsImplementation["google.golang.org/protobuf/types/known/timestamppb"] = true
					childTypeType = "*timestamppb.Timestamp"
					childTypeToAppend = "timestamppb.New(v)"
				}
				if childType.Type == schemas.TypeType_Enum {
					if childType.EnumHash == nil {
						return nil, nil, fmt.Errorf("enum \"%s\" not found", *childType.EnumHash)
					}

					schemaEnum := self.schema.Enums.Enums[*childType.EnumHash]
					enum, err := self.toEnum(schemaEnum)
					if err != nil {
						return nil, nil, err
					}

					childTypeType = "pb." + enum.Name
					childTypeToAppend = fmt.Sprintf("convert%sToPb(v)", enum.Name)
				}

				if childTypeType == "" {
					return nil, nil, fmt.Errorf("unable to parse \"%s\": grpc-client-go currently doesn't support lists of lists and lists of maps", childType.Name)
				}

				prepareList, err := self.templateManager.Parse("input-prop-list", &templates.InputPropListTemplInput{
					MethodName:           methodName,
					VarName:              varName,
					OriginalVariableName: fmt.Sprintf("i.%s", propType.Name),
					Type:                 childTypeType,
					ValueToAppend:        childTypeToAppend,
					Optional:             propType.Optional,
					ChildOptional:        childType.Optional,
				})
				if err != nil {
					return nil, nil, err
				}
				propsPrepare = append(propsPrepare, prepareList)

				value = varName
			}
		}
		if propType.Type == schemas.TypeType_Map {
			if propType.ChildTypesHashes == nil {
				return nil, nil, fmt.Errorf("ChildTypesHashes for \"%s\" not found", t.Name)
			}

			childBiggest := 0
			childTypes := []*schemas.Type{}
			for _, v := range propType.ChildTypesHashes {
				childPropType, ok := self.schema.Types.Types[v]
				if !ok {
					return nil, nil, fmt.Errorf("type \"%s\" not found", v)
				}

				if len(childPropType.Name) > childBiggest {
					childBiggest = len(childPropType.Name)
				}

				childTypes = append(childTypes, childPropType)
			}

			propsProps := []*templates.InputPropMapTemplProp{}

			for _, childChildType := range childTypes {
				var value string
				if childChildType.Type == schemas.TypeType_String ||
					childChildType.Type == schemas.TypeType_Int ||
					childChildType.Type == schemas.TypeType_Float ||
					childChildType.Type == schemas.TypeType_Bool {
					value = fmt.Sprintf("i.%s.%s", propType.Name, childChildType.Name)
				}
				if childChildType.Type == schemas.TypeType_Timestamp {
					if childChildType.Optional {
						return nil, nil, fmt.Errorf("grpc-client-go doesn't support optional map child timestamp properties")
					}

					self.importsContract["time"] = true
					value = fmt.Sprintf("i.%s.%s.AsTime()", propType.Name, childChildType.Name)
				}
				if childChildType.Type == schemas.TypeType_Enum {
					if childChildType.Optional {
						return nil, nil, fmt.Errorf("grpc-client-go doesn't support optional map child timestamp properties")
					}

					if childChildType.EnumHash == nil {
						return nil, nil, fmt.Errorf("enum \"%s\" not found", *childChildType.EnumHash)
					}

					schemaEnum := self.schema.Enums.Enums[*childChildType.EnumHash]
					enum, err := self.toEnum(schemaEnum)
					if err != nil {
						return nil, nil, err
					}

					value = fmt.Sprintf("convert%sToPb(i.%s.%s)", enum.Name, propType.Name, childChildType.Name)
				}

				if value == "" {
					return nil, nil, fmt.Errorf("unable to parse \"%s\": grpc-client-go currently doesn't support maps of lists and maps of maps", childChildType.Name)
				}

				propsProps = append(propsProps, &templates.InputPropMapTemplProp{
					Name:    childChildType.Name,
					Spacing: strings.Repeat(" ", childBiggest-len(childChildType.Name)),
					Value:   value,
				})
			}

			varName := formatter.PascalToCamel(propType.Name)

			prepareMap, err := self.templateManager.Parse("input-prop-map", &templates.InputPropMapTemplInput{
				VarName: varName,
				Type:    propType.Name,
				Props:   propsProps,
			})
			if err != nil {
				return nil, nil, err
			}
			propsPrepare = append(propsPrepare, prepareMap)

			value = varName
		}

		props = append(props, &templates.TemplMethodProp{
			Name:    propType.Name,
			Spacing: strings.Repeat(" ", biggest-len(propType.Name)),
			Value:   value,
		})
	}

	return props, propsPrepare, nil
}
