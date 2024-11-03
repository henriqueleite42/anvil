package grpc

import (
	"fmt"
	"strings"

	"github.com/henriqueleite42/anvil/language-helpers/golang/formatter"
	"github.com/henriqueleite42/anvil/language-helpers/golang/grpc/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/imports"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *goGrpcParser) goToProto(i *convertingInput) (*convertingValue, error) {
	if i == nil || i.input == nil {
		return nil, fmt.Errorf("GoToProto: input is required")
	}
	if i.input.Type == nil {
		return nil, fmt.Errorf("GoToProto: \"Type\" is required")
	}
	oi := i.input
	t := oi.Type

	name := t.Name
	if i.overwriteTypeName != nil {
		name = *i.overwriteTypeName
	}
	nameWithPrefix := name
	if i.prefixForVariableNaming != nil {
		nameWithPrefix = *i.prefixForVariableNaming + name
	}

	parsedType, err := self.goTypeParser.ParseType(t)
	if err != nil {
		return nil, err
	}

	golangType := parsedType.GetFullTypeName(oi.CurModuleImport.Alias)
	golangTypeName := parsedType.GetTypeName(oi.CurModuleImport.Alias)

	if doestNeedConversion(t.Type) {
		// Doesn't need conversion
		return &convertingValue{
			GolangType:     golangType,
			GolangTypeName: golangTypeName,
			ProtoType:      golangType,
			ProtoTypeName:  golangTypeName,
			Value:          oi.VarToConvert,
		}, nil
	}
	if t.Type == schemas.TypeType_Int ||
		t.Type == schemas.TypeType_Int8 ||
		t.Type == schemas.TypeType_Int16 {
		pbType := "int32"

		val := &convertingValue{
			GolangType:     golangType,
			GolangTypeName: golangTypeName,
			ProtoTypeName:  pbType,
		}

		if t.Optional {
			varName := formatter.PascalToCamel(nameWithPrefix)
			prepareOptional, err := self.templateManager.Parse("input-prop-optional", &templates.InputPropOptionalTemplInput{
				VarName:              varName,
				OriginalVariableName: oi.VarToConvert,
				Type:                 "*" + pbType,
				ValueToAssign:        fmt.Sprintf("%s(*%s)", pbType, oi.VarToConvert),
				NeedsPointer:         true,
			})
			if err != nil {
				return nil, err
			}

			val.ProtoType = "*" + pbType
			val.Value = varName
			val.Prepare = []string{prepareOptional}

			return val, nil
		} else {
			val.ProtoType = pbType
			val.Value = fmt.Sprintf("%s(%s)", pbType, oi.VarToConvert)

			return val, nil
		}
	}
	if t.Type == schemas.TypeType_Uint ||
		t.Type == schemas.TypeType_Uint8 ||
		t.Type == schemas.TypeType_Uint16 {
		pbType := "uint32"

		val := &convertingValue{
			GolangType:     golangType,
			GolangTypeName: golangTypeName,
			ProtoTypeName:  pbType,
		}

		if t.Optional {
			varName := formatter.PascalToCamel(nameWithPrefix)
			prepareOptional, err := self.templateManager.Parse("input-prop-optional", &templates.InputPropOptionalTemplInput{
				VarName:              varName,
				OriginalVariableName: oi.VarToConvert,
				Type:                 "*" + pbType,
				ValueToAssign:        fmt.Sprintf("%s(*%s)", pbType, oi.VarToConvert),
				NeedsPointer:         true,
			})
			if err != nil {
				return nil, err
			}

			val.ProtoType = "*" + pbType
			val.Value = varName
			val.Prepare = []string{prepareOptional}

			return val, nil
		} else {
			val.ProtoType = pbType
			val.Value = fmt.Sprintf("%s(%s)", pbType, oi.VarToConvert)

			return val, nil
		}
	}
	if t.Type == schemas.TypeType_Float {
		pbType := "float32"

		val := &convertingValue{
			GolangType:     golangType,
			GolangTypeName: golangTypeName,
			ProtoTypeName:  pbType,
		}

		if t.Optional {
			varName := formatter.PascalToCamel(nameWithPrefix)
			prepareOptional, err := self.templateManager.Parse("input-prop-optional", &templates.InputPropOptionalTemplInput{
				VarName:              varName,
				OriginalVariableName: oi.VarToConvert,
				Type:                 "*" + pbType,
				ValueToAssign:        fmt.Sprintf("%s(*%s)", pbType, oi.VarToConvert),
				NeedsPointer:         true,
			})
			if err != nil {
				return nil, err
			}

			val.ProtoType = "*" + pbType
			val.Value = varName
			val.Prepare = []string{prepareOptional}

			return val, nil
		} else {
			val.ProtoType = pbType
			val.Value = fmt.Sprintf("%s(%s)", pbType, oi.VarToConvert)

			return val, nil
		}
	}
	if t.Type == schemas.TypeType_Timestamp {
		importsManager := imports.NewImportsManager()
		importsManager.AddImport("google.golang.org/protobuf/types/known/timestamppb", nil)

		pbType := "timestamppb.Timestamp"

		val := &convertingValue{
			GolangType:     golangType,
			GolangTypeName: golangTypeName,
			ProtoType:      "*" + pbType,
			ProtoTypeName:  pbType,

			imports: importsManager,
		}

		if t.Optional {
			varName := formatter.PascalToCamel(nameWithPrefix)
			prepareOptional, err := self.templateManager.Parse("input-prop-optional", &templates.InputPropOptionalTemplInput{
				VarName:              varName,
				OriginalVariableName: oi.VarToConvert,
				Type:                 "*" + pbType,
				ValueToAssign:        fmt.Sprintf("timestamppb.New(*%s)", oi.VarToConvert),
			})
			if err != nil {
				return nil, err
			}

			val.Value = varName
			val.Prepare = []string{prepareOptional}

			return val, nil
		} else {
			val.Value = fmt.Sprintf("timestamppb.New(%s)", oi.VarToConvert)

			return val, nil
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

		importsManager := imports.NewImportsManager()
		importsManager.MergeImport(oi.PbModuleImport)

		pbType := oi.PbModuleImport.Alias + "." + enum.GolangName

		val := &convertingValue{
			GolangType:     golangType,
			GolangTypeName: golangTypeName,
			ProtoType:      "*" + pbType,
			ProtoTypeName:  pbType,
		}

		var firstLetter string
		var pkg string
		if self.enumConversionImport.Alias == oi.CurModuleImport.Alias {
			firstLetter = "c"
		} else {
			firstLetter = "C"
			pkg = self.enumConversionImport.Alias + "."
		}
		enumConvertFuncName := fmt.Sprintf("%s%sonvert%sToPb(*%s)", pkg, firstLetter, enum.GolangName, oi.VarToConvert)

		if t.Optional {
			varName := formatter.PascalToCamel(nameWithPrefix)
			prepareOptional, err := self.templateManager.Parse("input-prop-optional", &templates.InputPropOptionalTemplInput{
				VarName:              varName,
				OriginalVariableName: oi.VarToConvert,
				Type:                 "*" + pbType,
				ValueToAssign:        fmt.Sprintf("%s(*%s)", enumConvertFuncName, oi.VarToConvert),
				NeedsPointer:         true,
			})
			if err != nil {
				return nil, err
			}

			val.Value = varName
			val.Prepare = []string{prepareOptional}

			return val, nil
		} else {
			val.Value = fmt.Sprintf("%s(%s)", enumConvertFuncName, oi.VarToConvert)

			return val, nil
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

		if doestNeedConversion(childType.Type) {
			return self.goToProto(&convertingInput{
				input: &ConverterInput{
					CurModuleImport: oi.CurModuleImport,
					PbModuleImport:  oi.PbModuleImport,
					Type:            childType,
					VarToConvert:    oi.VarToConvert,
				},
			})
		}

		varName := formatter.PascalToCamel(nameWithPrefix)

		r, err := self.goToProto(&convertingInput{
			indentationLvl:          i.indentationLvl + 1,
			prefixForVariableNaming: &nameWithPrefix,
			input: &ConverterInput{
				CurModuleImport: oi.CurModuleImport,
				PbModuleImport:  oi.PbModuleImport,
				VarToConvert:    "v",
				Type:            childType,
			},
		})
		if err != nil {
			return nil, err
		}

		prepareList, err := self.templateManager.Parse("input-prop-list", &templates.InputPropListTemplInput{
			VarName:              varName,
			OriginalVariableName: oi.VarToConvert,
			Type:                 r.ProtoType,
			ValueToAppend:        r.Value,
			ChildOptional:        childType.Optional,
			Prepare:              r.Prepare,
		})
		if err != nil {
			return nil, err
		}

		return &convertingValue{
			GolangType:     golangType,
			GolangTypeName: golangTypeName,
			ProtoType:      "[]" + r.ProtoType,
			ProtoTypeName:  "[]" + r.ProtoType,
			Value:          varName,
			Prepare:        []string{prepareList},

			imports: r.imports,
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

		varName := formatter.PascalToCamel(nameWithPrefix)

		props := make([]*templates.InputPropMapTemplProp, 0, len(t.ChildTypes))
		var prepare []string = nil

		importsManager := imports.NewImportsManager()
		importsManager.MergeImports(parsedType.GetImportsUnorganized())

		for _, v := range t.ChildTypes {
			propType, ok := self.schema.Types.Types[v.TypeHash]
			if !ok {
				return nil, fmt.Errorf("type \"%s\" not found", v.TypeHash)
			}

			propNameWithPrefix := fmt.Sprintf("%s.%s", oi.VarToConvert, *v.PropName)

			if doestNeedConversion(propType.Type) {
				props = append(props, &templates.InputPropMapTemplProp{
					Name:    *v.PropName,
					Spacing: strings.Repeat(" ", biggestName-len(*v.PropName)),
					Value:   propNameWithPrefix,
				})
				continue
			}

			overwriteTypeName := name + *v.PropName
			r, err := self.goToProto(&convertingInput{
				indentationLvl:          i.indentationLvl + 1,
				overwriteTypeName:       &overwriteTypeName,
				prefixForVariableNaming: &nameWithPrefix,

				input: &ConverterInput{
					CurModuleImport: oi.CurModuleImport,
					PbModuleImport:  oi.PbModuleImport,
					Type:            propType,
					VarToConvert:    propNameWithPrefix,
				},
			})
			if err != nil {
				return nil, err
			}

			if r.imports != nil {
				importsManager.MergeImports(r.imports.GetImportsUnorganized())
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

		pbType, err := GetProtoTypeName(t)
		if err != nil {
			return nil, err
		}

		prepareMap, err := self.templateManager.Parse("input-prop-map", &templates.InputPropMapTemplInput{
			VarName:              varName,
			OriginalVariableName: oi.VarToConvert,
			Type:                 pbType,
			Prepare:              prepare,
			IndentationLvl:       i.indentationLvl + 1,
			Props:                props,
			TypePkg:              &oi.PbModuleImport.Alias,
		})
		if err != nil {
			return nil, err
		}

		pbTypeWithPkg := oi.PbModuleImport.Alias + "." + pbType

		var mapImports imports.ImportsManager = nil
		if importsManager.GetImportsLen() != 0 {
			mapImports = importsManager
		}

		return &convertingValue{
			GolangType:     golangType,
			GolangTypeName: golangTypeName,
			ProtoType:      "*" + pbTypeWithPkg,
			ProtoTypeName:  pbTypeWithPkg,
			Value:          varName,
			Prepare:        []string{prepareMap},

			imports: mapImports,
		}, nil
	}

	return nil, fmt.Errorf("invalid type \"%s\"", t.Type)
}

func (self *goGrpcParser) GoToProto(i *ConverterInput) (*ConvertedValue, error) {
	result, err := self.goToProto(&convertingInput{
		input:          i,
		indentationLvl: 0,
	})
	if err != nil {
		return nil, err
	}

	var importsUnorganized []*imports.Import = nil
	if result.imports != nil {
		importsUnorganized = result.imports.GetImportsUnorganized()
	}

	return &ConvertedValue{
		GolangType:     result.GolangType,
		GolangTypeName: result.GolangTypeName,
		ProtoType:      result.ProtoType,
		ProtoTypeName:  result.ProtoTypeName,

		Value:              result.Value,
		Prepare:            result.Prepare,
		ImportsUnorganized: importsUnorganized,
	}, nil
}
