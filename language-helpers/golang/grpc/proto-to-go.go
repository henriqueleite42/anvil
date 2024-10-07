package grpc

import (
	"fmt"
	"strings"

	"github.com/henriqueleite42/anvil/language-helpers/golang/formatter"
	"github.com/henriqueleite42/anvil/language-helpers/golang/grpc/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *goGrpcParser) ProtoToGo(i *ProtoToGoInput) (*Type, error) {
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

	_, err := self.goTypeParser.ParseType(t)
	if err != nil {
		return nil, err
	}

	result := &Type{
		Name:         t.Name,
		Props:        make([]*Prop, 0, len(t.ChildTypes)),
		PropsPrepare: []string{},
	}

	biggest := 0
	for k, v := range t.ChildTypes {
		if v.PropName == nil {
			return nil, fmt.Errorf("ChildType \"%s.%d\" must have a PropName", t.Name, k)
		}

		if len(*v.PropName) > biggest {
			biggest = len(*v.PropName)
		}
	}

	for _, v := range t.ChildTypes {
		propType, ok := self.schema.Types.Types[v.TypeHash]
		if !ok {
			return nil, fmt.Errorf("type \"%s\" not found", v.TypeHash)
		}

		var value string

		propNameWithPrefix := *v.PropName
		if i.VariableName != "" {
			propNameWithPrefix = fmt.Sprintf("%s.%s", i.VariableName, *v.PropName)
		}

		if propType.Type == schemas.TypeType_String ||
			propType.Type == schemas.TypeType_Int ||
			propType.Type == schemas.TypeType_Float ||
			propType.Type == schemas.TypeType_Bool {
			value = propNameWithPrefix
		}
		if propType.Type == schemas.TypeType_Timestamp {
			if propType.Optional {
				self.goTypeParser.AddImport("time")
				varName := formatter.PascalToCamel(i.PrefixForVariableNaming + *v.PropName)
				prepareList, err := self.templateManager.Parse("input-prop-optional", &templates.InputPropOptionalTemplInput{
					VarName:              varName,
					OriginalVariableName: propNameWithPrefix,
					Type:                 "*time.Time",
					ValueToAssign:        fmt.Sprintf("%s.AsTime()", propNameWithPrefix),
					NeedsPointer:         true,
				})
				if err != nil {
					return nil, err
				}
				result.PropsPrepare = append(result.PropsPrepare, prepareList)
				value = varName
			} else {
				value = fmt.Sprintf("%s.AsTime()", propNameWithPrefix)
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
				eType := "*" + enum.GetFullEnumName(i.CurPkg)

				varName := formatter.PascalToCamel(i.PrefixForVariableNaming + *v.PropName)
				prepareList, err := self.templateManager.Parse("input-prop-optional", &templates.InputPropOptionalTemplInput{
					VarName:              varName,
					OriginalVariableName: propNameWithPrefix,
					Type:                 eType,
					ValueToAssign:        fmt.Sprintf("convertPbTo%s(*%s)", enum.GolangName, propNameWithPrefix),
					NeedsPointer:         true,
				})
				if err != nil {
					return nil, err
				}
				result.PropsPrepare = append(result.PropsPrepare, prepareList)
				value = varName
			} else {
				value = fmt.Sprintf("convertPbTo%s(%s)", enum.GolangName, propNameWithPrefix)
			}
		}
		if propType.Type == schemas.TypeType_List {
			if propType.Optional {
				return nil, fmt.Errorf("unable to parse \"%s\": language-helper golang grpc currently doesn't support optional lists", *v.PropName)
			}
			if propType.ChildTypes == nil {
				return nil, fmt.Errorf("ChildTypes for \"%s\" not found", t.Name)
			}
			if len(propType.ChildTypes) != 1 {
				return nil, fmt.Errorf("ChildTypes for \"%s\" must have exactly one item", t.Name)
			}

			childType, ok := self.schema.Types.Types[propType.ChildTypes[0].TypeHash]
			if !ok {
				return nil, fmt.Errorf("type \"%s\" not found", propType.ChildTypes[0].TypeHash)
			}

			if childType.Type == schemas.TypeType_String ||
				childType.Type == schemas.TypeType_Int ||
				childType.Type == schemas.TypeType_Float ||
				childType.Type == schemas.TypeType_Bool {
				value = propNameWithPrefix
			} else {
				var childTypeType string
				var childTypeToAppend string
				var prepare []string = nil

				varNamePascal := i.PrefixForVariableNaming + *v.PropName
				varName := formatter.PascalToCamel(varNamePascal)

				if childType.Type == schemas.TypeType_Timestamp {
					childTypeType = "time.Time"
					childTypeToAppend = "v.AsTime()"
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

					childTypeType = enum.GetFullEnumName(i.CurPkg)
					childTypeToAppend = fmt.Sprintf("convertPbTo%s(v)", enum.GolangName)
				}
				if childType.Type == schemas.TypeType_Map {
					r, err := self.ProtoToGo(&ProtoToGoInput{
						indentationLvl:          i.indentationLvl + 1,
						Type:                    childType,
						MethodName:              i.MethodName,
						VariableName:            "v",
						PrefixForVariableNaming: varNamePascal,
						HasOutput:               i.HasOutput,
						CurPkg:                  i.CurPkg,
					})
					if err != nil {
						return nil, err
					}

					childTypeParsed, err := self.goTypeParser.ParseType(childType)
					if err != nil {
						return nil, err
					}

					propsProps := make([]*templates.InputPropMapTemplProp, len(r.Props), len(r.Props))
					for k, v := range r.Props {
						propsProps[k] = &templates.InputPropMapTemplProp{
							Name:    v.Name,
							Spacing: v.Spacing,
							Value:   v.Value,
						}
					}

					childVarName := formatter.PascalToCamel(varNamePascal + "Item")
					var childTypePkg *string
					if childTypeParsed.GolangPkg != nil && *childTypeParsed.GolangPkg != i.CurPkg {
						childTypePkg = childTypeParsed.GolangPkg
					}

					prepareMap, err := self.templateManager.Parse("input-prop-map", &templates.InputPropMapTemplInput{
						MethodName:           i.MethodName,
						Optional:             t.Optional,
						HasOutput:            i.HasOutput,
						OriginalVariableName: "v",
						TypePkg:              childTypePkg,
						VarName:              childVarName,
						Props:                propsProps,
						Type:                 childTypeParsed.GolangType,
						Prepare:              r.PropsPrepare,
						IndentationLvl:       i.indentationLvl,
					})
					if err != nil {
						return nil, err
					}

					prepare = []string{prepareMap}
					if childTypePkg != nil {
						childTypeType = fmt.Sprintf("*%s.%s", *childTypePkg, childTypeParsed.GolangType)
					} else {
						childTypeType = "*" + childTypeParsed.GolangType
					}
					childTypeToAppend = childVarName
				}

				if childTypeType == "" {
					return nil, fmt.Errorf("unable to parse \"%s\": language-helper golang grpc currently doesn't support lists of lists and lists of maps", childType.Name)
				}

				prepareList, err := self.templateManager.Parse("input-prop-list", &templates.InputPropListTemplInput{
					MethodName:           i.MethodName,
					VarName:              varName,
					OriginalVariableName: propNameWithPrefix,
					Type:                 childTypeType,
					ValueToAppend:        childTypeToAppend,
					Optional:             propType.Optional,
					ChildOptional:        childType.Optional,
					HasOutput:            i.HasOutput,
					Prepare:              prepare,
				})
				if err != nil {
					return nil, err
				}
				result.PropsPrepare = append(result.PropsPrepare, prepareList)

				value = varName
			}
		}
		if propType.Type == schemas.TypeType_Map {
			if propType.ChildTypes == nil {
				return nil, fmt.Errorf("ChildTypes for \"%s\" not found", t.Name)
			}

			childBiggest := 0
			for k, v := range propType.ChildTypes {
				if v.PropName == nil {
					return nil, fmt.Errorf("ChildType \"%s.%d\" must have a PropName", t.Name, k)
				}

				if len(*v.PropName) > childBiggest {
					childBiggest = len(*v.PropName)
				}
			}

			varNamePascal := i.PrefixForVariableNaming + *v.PropName
			varName := formatter.PascalToCamel(varNamePascal)
			propsProps := []*templates.InputPropMapTemplProp{}
			childPropsPrepare := []string{}
			for _, vv := range propType.ChildTypes {
				childPropType, ok := self.schema.Types.Types[vv.TypeHash]
				if !ok {
					return nil, fmt.Errorf("type \"%s\" not found", vv.TypeHash)
				}

				var value string
				if childPropType.Type == schemas.TypeType_String ||
					childPropType.Type == schemas.TypeType_Int ||
					childPropType.Type == schemas.TypeType_Float ||
					childPropType.Type == schemas.TypeType_Bool {
					value = fmt.Sprintf("%s.%s", propNameWithPrefix, *vv.PropName)
				}
				if childPropType.Type == schemas.TypeType_Timestamp {
					if childPropType.Optional {
						return nil, fmt.Errorf("language-helper golang grpc ProtoToGo doesn't support optional map child timestamp properties")
					}

					value = fmt.Sprintf("%s.%s.AsTime()", propNameWithPrefix, *vv.PropName)
				}
				if childPropType.Type == schemas.TypeType_Enum {
					if childPropType.Optional {
						return nil, fmt.Errorf("language-helper golang grpc ProtoToGo doesn't support optional map child enum properties")
					}

					if childPropType.EnumHash == nil {
						return nil, fmt.Errorf("enum \"%s\" not found", *childPropType.EnumHash)
					}

					schemaEnum := self.schema.Enums.Enums[*childPropType.EnumHash]
					enum, err := self.goTypeParser.ParseEnum(schemaEnum)
					if err != nil {
						return nil, err
					}

					value = fmt.Sprintf("convertPbTo%s(%s.%s)", enum.GolangName, propNameWithPrefix, *vv.PropName)
				}
				if childPropType.Type == schemas.TypeType_Map {
					childVarNamePascal := varNamePascal + *vv.PropName
					childVarName := formatter.PascalToCamel(childVarNamePascal)
					childPropWithPrefix := fmt.Sprintf("%s.%s", propNameWithPrefix, *vv.PropName)

					r, err := self.ProtoToGo(&ProtoToGoInput{
						Type:                    childPropType,
						MethodName:              i.MethodName,
						VariableName:            childPropWithPrefix,
						PrefixForVariableNaming: childVarNamePascal,
						HasOutput:               i.HasOutput,
						CurPkg:                  i.CurPkg,

						indentationLvl: i.indentationLvl + 1,
					})
					if err != nil {
						return nil, err
					}

					childProps := make([]*templates.InputPropMapTemplProp, len(r.Props), len(r.Props))
					for k, v := range r.Props {
						childProps[k] = &templates.InputPropMapTemplProp{
							Name:    v.Name,
							Spacing: v.Spacing,
							Value:   v.Value,
						}
					}

					childPropParsed, err := self.goTypeParser.ParseType(childPropType)
					if err != nil {
						return nil, err
					}

					var typePkg *string
					if childPropParsed.GolangPkg != nil && *childPropParsed.GolangPkg != i.CurPkg {
						typePkg = childPropParsed.GolangPkg
					}

					prepareChildMap, err := self.templateManager.Parse("input-prop-map", &templates.InputPropMapTemplInput{
						MethodName:           i.MethodName,
						Optional:             childPropType.Optional,
						HasOutput:            i.HasOutput,
						OriginalVariableName: childPropWithPrefix,
						TypePkg:              typePkg,
						VarName:              childVarName,
						Type:                 r.Name,
						Props:                childProps,
						Prepare:              r.PropsPrepare,
						IndentationLvl:       i.indentationLvl,
					})
					if err != nil {
						return nil, err
					}

					childPropsPrepare = append(childPropsPrepare, prepareChildMap)
					value = childVarName
				}

				if value == "" {
					return nil, fmt.Errorf("unable to parse \"%s\": language-helper golang grpc something went wrong", *vv.PropName)
				}

				propsProps = append(propsProps, &templates.InputPropMapTemplProp{
					Name:    *vv.PropName,
					Spacing: strings.Repeat(" ", childBiggest-len(*vv.PropName)),
					Value:   value,
				})
			}

			propTypeParsed, err := self.goTypeParser.ParseType(propType)
			if err != nil {
				return nil, err
			}

			var typePkg *string
			if propTypeParsed.GolangPkg != nil && *propTypeParsed.GolangPkg != i.CurPkg {
				typePkg = propTypeParsed.GolangPkg
			}

			prepareMap, err := self.templateManager.Parse("input-prop-map", &templates.InputPropMapTemplInput{
				MethodName:           i.MethodName,
				Optional:             t.Optional,
				HasOutput:            i.HasOutput,
				OriginalVariableName: propNameWithPrefix,
				TypePkg:              typePkg,
				VarName:              varName,
				Props:                propsProps,
				Type:                 propTypeParsed.GolangType,
				Prepare:              childPropsPrepare,
				IndentationLvl:       i.indentationLvl,
			})
			if err != nil {
				return nil, err
			}
			result.PropsPrepare = append(result.PropsPrepare, prepareMap)

			value = varName
		}

		result.Props = append(result.Props, &Prop{
			Name:    *v.PropName,
			Spacing: strings.Repeat(" ", biggest-len(*v.PropName)),
			Value:   value,
		})
	}

	return result, nil
}
