package grpc

import (
	"fmt"
	"strings"

	"github.com/henriqueleite42/anvil/language-helpers/golang/formatter"
	"github.com/henriqueleite42/anvil/language-helpers/golang/grpc/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *goGrpcParser) GoToProto(i *GoToProtoInput) (*Type, error) {
	if i == nil {
		return nil, fmt.Errorf("GoToProto: input is required")
	}
	if i.Type == nil {
		return nil, fmt.Errorf("GoToProto: \"Type\" is required")
	}
	t := i.Type
	if t.Type != schemas.TypeType_Map {
		return nil, fmt.Errorf("inputs for grpc must be Map Type")
	}

	_, err := self.goTypeParser.ParseType(t, nil)
	if err != nil {
		return nil, err
	}

	result := &Type{
		Name:         t.Name,
		Props:        make([]*Prop, 0, len(t.ChildTypesHashes)),
		PropsPrepare: []string{},
	}

	biggest := 0
	amountOfOptional := 0
	types := []*schemas.Type{}
	for _, v := range t.ChildTypesHashes {
		propType, ok := self.schema.Types.Types[v]
		if !ok {
			return nil, fmt.Errorf("type \"%s\" not found", v)
		}

		types = append(types, propType)

		if len(propType.Name) > biggest {
			biggest = len(propType.Name)
		}

		if propType.Optional {
			amountOfOptional++
		}
	}

	for _, propType := range types {
		var value string

		propNameWithPrefix := propType.Name
		if i.Prefix != "" {
			propNameWithPrefix = fmt.Sprintf("%s.%s", i.Prefix, propType.Name)
		}

		if propType.Type == schemas.TypeType_String ||
			propType.Type == schemas.TypeType_Int ||
			propType.Type == schemas.TypeType_Float ||
			propType.Type == schemas.TypeType_Bool {
			value = "i." + propType.Name
		}
		if propType.Type == schemas.TypeType_Timestamp {
			self.goTypeParser.AddImport("google.golang.org/protobuf/types/known/timestamppb")

			if propType.Optional {
				varName := formatter.PascalToCamel(propType.Name)
				prepareList, err := self.templateManager.Parse("input-prop-optional", &templates.InputPropOptionalTemplInput{
					VarName:              varName,
					OriginalVariableName: propNameWithPrefix,
					Type:                 "*timestamppb.Timestamp",
					ValueToAssign:        fmt.Sprintf("timestamppb.New(*%s)", propNameWithPrefix),
				})
				if err != nil {
					return nil, err
				}
				result.PropsPrepare = append(result.PropsPrepare, prepareList)
				value = varName
			} else {
				value = fmt.Sprintf("timestamppb.New(%s)", propNameWithPrefix)
			}
		}
		if propType.Type == schemas.TypeType_Enum {
			if propType.EnumHash == nil {
				return nil, fmt.Errorf("enum \"%s\" not found", *propType.EnumHash)
			}

			schemaEnum := self.schema.Enums.Enums[*propType.EnumHash]
			enum, err := self.goTypeParser.ParseEnum(schemaEnum)
			if err != nil {
				return nil, err
			}

			if propType.Optional {
				varName := formatter.PascalToCamel(propType.Name)
				prepareList, err := self.templateManager.Parse("input-prop-optional", &templates.InputPropOptionalTemplInput{
					VarName:              varName,
					OriginalVariableName: propNameWithPrefix,
					Type:                 "*pb." + enum.GolangName,
					ValueToAssign:        fmt.Sprintf("convert%sToPb(*%s)", enum.GolangName, propNameWithPrefix),
					NeedsPointer:         true,
				})
				if err != nil {
					return nil, err
				}
				result.PropsPrepare = append(result.PropsPrepare, prepareList)
				value = varName
			} else {
				value = fmt.Sprintf("convert%sToPb(%s)", enum.GolangName, propNameWithPrefix)
			}
		}
		if propType.Type == schemas.TypeType_List {
			if propType.Optional {
				return nil, fmt.Errorf("unable to parse \"%s\": grpc-client-go currently doesn't support optional lists", propType.Name)
			}
			if propType.ChildTypesHashes == nil {
				return nil, fmt.Errorf("ChildTypesHashes for \"%s\" not found", t.Name)
			}
			if len(propType.ChildTypesHashes) != 1 {
				return nil, fmt.Errorf("ChildTypesHashes for \"%s\" must have exactly one item", t.Name)
			}

			childType, ok := self.schema.Types.Types[propType.ChildTypesHashes[0]]
			if !ok {
				return nil, fmt.Errorf("type \"%s\" not found", propType.ChildTypesHashes[0])
			}

			varName := formatter.PascalToCamel(propType.Name)

			if childType.Type == schemas.TypeType_String ||
				childType.Type == schemas.TypeType_Int ||
				childType.Type == schemas.TypeType_Float ||
				childType.Type == schemas.TypeType_Bool {
				value = propNameWithPrefix
			} else {
				var childTypeType string
				var childTypeToAppend string
				if childType.Type == schemas.TypeType_Timestamp {
					self.goTypeParser.AddImport("google.golang.org/protobuf/types/known/timestamppb")
					childTypeType = "*timestamppb.Timestamp"
					childTypeToAppend = "timestamppb.New(v)"
				}
				if childType.Type == schemas.TypeType_Enum {
					if childType.EnumHash == nil {
						return nil, fmt.Errorf("enum \"%s\" not found", *childType.EnumHash)
					}

					schemaEnum := self.schema.Enums.Enums[*childType.EnumHash]
					enum, err := self.goTypeParser.ParseEnum(schemaEnum)
					if err != nil {
						return nil, err
					}

					childTypeType = "pb." + enum.GolangName
					childTypeToAppend = fmt.Sprintf("convert%sToPb(v)", enum.GolangName)
				}

				if childTypeType == "" {
					return nil, fmt.Errorf("unable to parse \"%s\": grpc-client-go currently doesn't support lists of lists and lists of maps", childType.Name)
				}

				prepareList, err := self.templateManager.Parse("input-prop-list", &templates.InputPropListTemplInput{
					MethodName:           i.MethodName,
					VarName:              varName,
					OriginalVariableName: propNameWithPrefix,
					Type:                 childTypeType,
					ValueToAppend:        childTypeToAppend,
					Optional:             propType.Optional,
					ChildOptional:        childType.Optional,
				})
				if err != nil {
					return nil, err
				}
				result.PropsPrepare = append(result.PropsPrepare, prepareList)

				value = varName
			}
		}
		if propType.Type == schemas.TypeType_Map {
			if propType.ChildTypesHashes == nil {
				return nil, fmt.Errorf("ChildTypesHashes for \"%s\" not found", t.Name)
			}

			childBiggest := 0
			childTypes := []*schemas.Type{}
			for _, v := range propType.ChildTypesHashes {
				childPropType, ok := self.schema.Types.Types[v]
				if !ok {
					return nil, fmt.Errorf("type \"%s\" not found", v)
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
					value = fmt.Sprintf("%s.%s", propNameWithPrefix, childChildType.Name)
				}
				if childChildType.Type == schemas.TypeType_Timestamp {
					if childChildType.Optional {
						return nil, fmt.Errorf("grpc-client-go doesn't support optional map child timestamp properties")
					}

					value = fmt.Sprintf("%s.%s.AsTime()", propNameWithPrefix, childChildType.Name)
				}
				if childChildType.Type == schemas.TypeType_Enum {
					if childChildType.Optional {
						return nil, fmt.Errorf("grpc-client-go doesn't support optional map child timestamp properties")
					}

					if childChildType.EnumHash == nil {
						return nil, fmt.Errorf("enum \"%s\" not found", *childChildType.EnumHash)
					}

					schemaEnum := self.schema.Enums.Enums[*childChildType.EnumHash]
					enum, err := self.goTypeParser.ParseEnum(schemaEnum)
					if err != nil {
						return nil, err
					}

					value = fmt.Sprintf("convert%sToPb(%s.%s)", enum.GolangName, propNameWithPrefix, childChildType.Name)
				}

				if value == "" {
					return nil, fmt.Errorf("unable to parse \"%s\": grpc-client-go currently doesn't support maps of lists and maps of maps", childChildType.Name)
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
				return nil, err
			}
			result.PropsPrepare = append(result.PropsPrepare, prepareMap)

			value = varName
		}

		result.Props = append(result.Props, &Prop{
			Name:    propType.Name,
			Spacing: strings.Repeat(" ", biggest-len(propType.Name)),
			Value:   value,
		})
	}

	return result, nil
}
