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
		return nil, fmt.Errorf("ProtoToGo: input is required")
	}
	if i.Type == nil {
		return nil, fmt.Errorf("ProtoToGo: \"Type\" is required")
	}
	t := i.Type

	name := i.TypeName
	if name == "" {
		name = t.Name
	}

	parsedType, err := self.goTypeParser.ParseType(t)
	if err != nil {
		return nil, err
	}

	golangType := parsedType.GetFullTypeName(i.CurPkg)
	golangTypeName := parsedType.GetTypeName(i.CurPkg)

	if t.Type == schemas.TypeType_String ||
		t.Type == schemas.TypeType_Int ||
		t.Type == schemas.TypeType_Float ||
		t.Type == schemas.TypeType_Bool {
		return &Type{
			GolangType:     golangType,
			GolangTypeName: golangTypeName,
			ProtoType:      golangType,
			ProtoTypeName:  golangTypeName,
			Value:          i.VariableToAccessTheValue,
		}, nil
	}
	if t.Type == schemas.TypeType_Timestamp {
		self.goTypeParser.AddImport("time")

		pbType := "timestamppb.Timestamp"

		if t.Optional {
			varName := formatter.PascalToCamel(i.PrefixForVariableNaming + name)
			prepareOptional, err := self.templateManager.Parse("input-prop-optional", &templates.InputPropOptionalTemplInput{
				VarName:              varName,
				OriginalVariableName: i.VariableToAccessTheValue,
				Type:                 golangType,
				ValueToAssign:        fmt.Sprintf("%s.AsTime()", i.VariableToAccessTheValue),
				NeedsPointer:         true,
			})
			if err != nil {
				return nil, err
			}

			return &Type{
				GolangType:     golangType,
				GolangTypeName: golangTypeName,
				ProtoType:      "*" + pbType,
				ProtoTypeName:  pbType,
				Value:          varName,
				Prepare:        []string{prepareOptional},
			}, nil
		} else {
			return &Type{
				GolangType:     golangType,
				GolangTypeName: golangTypeName,
				ProtoType:      "*" + pbType,
				ProtoTypeName:  pbType,
				Value:          fmt.Sprintf("%s.AsTime()", i.VariableToAccessTheValue),
			}, nil
		}
	}
	if t.Type == schemas.TypeType_Enum {
		if t.EnumHash == nil {
			return nil, fmt.Errorf("enum for type \"%s\" not found", name)
		}

		schemaEnum, ok := self.schema.Enums.Enums[*t.EnumHash]
		if !ok {
			return nil, fmt.Errorf("enum \"%s\" not found", *t.EnumHash)
		}
		enum, err := self.goTypeParser.ParseEnum(schemaEnum)
		if err != nil {
			return nil, err
		}

		pbType := "pb." + enum.GolangName

		if t.Optional {
			varName := formatter.PascalToCamel(i.PrefixForVariableNaming + name)
			prepareOptional, err := self.templateManager.Parse("input-prop-optional", &templates.InputPropOptionalTemplInput{
				VarName:              varName,
				OriginalVariableName: i.VariableToAccessTheValue,
				Type:                 parsedType.GetFullTypeName(i.CurPkg),
				ValueToAssign:        fmt.Sprintf("convertPbTo%s(*%s)", enum.GolangName, i.VariableToAccessTheValue),
				NeedsPointer:         true,
			})
			if err != nil {
				return nil, err
			}

			return &Type{
				GolangType:     golangType,
				GolangTypeName: golangTypeName,
				ProtoType:      "*" + pbType,
				ProtoTypeName:  pbType,
				Prepare:        []string{prepareOptional},
				Value:          varName,
			}, nil
		} else {
			return &Type{
				GolangType:     golangType,
				GolangTypeName: golangTypeName,
				ProtoType:      pbType,
				ProtoTypeName:  pbType,
				Value:          fmt.Sprintf("convertPbTo%s(%s)", enum.GolangName, i.VariableToAccessTheValue),
			}, nil
		}
	}
	if t.Type == schemas.TypeType_List {
		if t.ChildTypes == nil {
			return nil, fmt.Errorf("ChildTypes for \"%s\" not found", name)
		}
		if len(t.ChildTypes) != 1 {
			return nil, fmt.Errorf("ChildTypes for \"%s\" must have exactly one item", name)
		}

		childType, ok := self.schema.Types.Types[t.ChildTypes[0].TypeHash]
		if !ok {
			return nil, fmt.Errorf("type \"%s\" not found", t.ChildTypes[0].TypeHash)
		}

		if childType.Type == schemas.TypeType_String ||
			childType.Type == schemas.TypeType_Int ||
			childType.Type == schemas.TypeType_Float ||
			childType.Type == schemas.TypeType_Bool {
			return &Type{
				GolangType:     golangType,
				GolangTypeName: golangTypeName,
				ProtoType:      golangType,
				ProtoTypeName:  golangTypeName,
				Value:          i.VariableToAccessTheValue,
			}, nil
		}

		varNamePascal := i.PrefixForVariableNaming + name
		varName := formatter.PascalToCamel(varNamePascal)

		r, err := self.ProtoToGo(&ProtoToGoInput{
			indentationLvl: i.indentationLvl + 1,
			MethodName:     i.MethodName,
			HasOutput:      i.HasOutput,
			CurPkg:         i.CurPkg,

			Type:                     childType,
			TypeName:                 name + "Item",
			VariableToAccessTheValue: "v",
			PrefixForVariableNaming:  varNamePascal,
		})
		if err != nil {
			return nil, err
		}

		prepareList, err := self.templateManager.Parse("input-prop-list", &templates.InputPropListTemplInput{
			MethodName:           i.MethodName,
			VarName:              varName,
			OriginalVariableName: i.VariableToAccessTheValue,
			Type:                 r.GolangType,
			ValueToAppend:        r.Value,
			Optional:             t.Optional,
			ChildOptional:        childType.Optional,
			HasOutput:            i.HasOutput,
			Prepare:              r.Prepare,
		})
		if err != nil {
			return nil, err
		}

		return &Type{
			GolangType:     golangType,
			GolangTypeName: golangTypeName,
			ProtoType:      "[]" + r.ProtoType,
			ProtoTypeName:  "[]" + r.ProtoType,
			Value:          varName,
			Prepare:        []string{prepareList},
		}, nil
	}
	if t.Type == schemas.TypeType_Map {
		if t.ChildTypes == nil {
			return nil, fmt.Errorf("ChildTypes for \"%s\" not found", name)
		}

		biggestName := 0
		for k, v := range t.ChildTypes {
			if v.PropName == nil {
				return nil, fmt.Errorf("ChildType \"%s.%d\" must have a PropName", t.Name, k)
			}

			if len(*v.PropName) > biggestName {
				biggestName = len(*v.PropName)
			}
		}

		varNamePascal := i.PrefixForVariableNaming + name
		varName := formatter.PascalToCamel(varNamePascal)

		props := make([]*templates.InputPropMapTemplProp, 0, len(t.ChildTypes))
		var prepare []string = nil

		for _, v := range t.ChildTypes {
			propType, ok := self.schema.Types.Types[v.TypeHash]
			if !ok {
				return nil, fmt.Errorf("type \"%s\" not found", v.TypeHash)
			}

			propNameWithPrefix := *v.PropName
			if i.VariableToAccessTheValue != "" {
				propNameWithPrefix = fmt.Sprintf("%s.%s", i.VariableToAccessTheValue, *v.PropName)
			}

			if propType.Type == schemas.TypeType_String ||
				propType.Type == schemas.TypeType_Int ||
				propType.Type == schemas.TypeType_Float ||
				propType.Type == schemas.TypeType_Bool {
				props = append(props, &templates.InputPropMapTemplProp{
					Name:    *v.PropName,
					Spacing: strings.Repeat(" ", biggestName-len(*v.PropName)),
					Value:   propNameWithPrefix,
				})
				continue
			}

			r, err := self.ProtoToGo(&ProtoToGoInput{
				indentationLvl: i.indentationLvl + 1,
				MethodName:     i.MethodName,
				HasOutput:      i.HasOutput,
				CurPkg:         i.CurPkg,

				Type:                     propType,
				TypeName:                 name + *v.PropName,
				VariableToAccessTheValue: propNameWithPrefix,
				PrefixForVariableNaming:  varNamePascal,
			})
			if err != nil {
				return nil, err
			}
			props = append(props, &templates.InputPropMapTemplProp{
				Name:    *v.PropName,
				Spacing: strings.Repeat(" ", biggestName-len(*v.PropName)),
				Value:   r.Value,
			})
			if prepare == nil {
				prepare = []string{}
			}
			prepare = append(prepare, r.Prepare...)
		}

		childTypeParsed, err := self.goTypeParser.ParseType(t)
		if err != nil {
			return nil, err
		}

		var typePkg *string
		if childTypeParsed.GolangPkg != nil && *childTypeParsed.GolangPkg != i.CurPkg {
			typePkg = childTypeParsed.GolangPkg
		}

		prepareMap, err := self.templateManager.Parse("input-prop-map", &templates.InputPropMapTemplInput{
			MethodName:           i.MethodName,
			VarName:              varName,
			OriginalVariableName: i.VariableToAccessTheValue,
			Type:                 parsedType.GolangType,
			Optional:             t.Optional,
			HasOutput:            i.HasOutput,
			Prepare:              prepare,
			IndentationLvl:       i.indentationLvl + 1,
			Props:                props,
			TypePkg:              typePkg,
		})
		if err != nil {
			return nil, err
		}

		pbType := "pb." + t.Name

		return &Type{
			GolangType:     golangType,
			GolangTypeName: golangTypeName,
			ProtoType:      "*" + pbType,
			ProtoTypeName:  pbType,
			Value:          varName,
			Prepare:        []string{prepareMap},
		}, nil
	}

	return nil, fmt.Errorf("invalid type \"%s\"", t.Type)
}
