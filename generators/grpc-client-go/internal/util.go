package internal

import (
	"fmt"
	"strings"

	"github.com/henriqueleite42/anvil/generators/grpc-client-go/internal/templates"
	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

func typeToTemplType(curPkg string, t *types_parser.Type) *templates.TemplType {
	var mapProps []*templates.TemplTypeMapProp = nil
	if t.MapProps != nil {
		biggestName := 0
		biggestType := 0
		for _, v := range t.MapProps {
			lenName := len(v.Name)
			if lenName > biggestName {
				biggestName = lenName
			}
			tType := v.Type.GetFullTypeName(curPkg)
			lenType := len(tType)
			if lenType > biggestType {
				biggestType = lenType
			}
		}

		mapProps = make([]*templates.TemplTypeMapProp, 0, len(t.MapProps))

		for _, v := range t.MapProps {
			tType := v.Type.GetFullTypeName(curPkg)
			mapProps = append(mapProps, &templates.TemplTypeMapProp{
				Name:     v.Name,
				Spacing1: strings.Repeat(" ", biggestName-len(v.Name)),
				Type:     tType,
				Spacing2: strings.Repeat(" ", biggestType-len(tType)),
				Tags:     v.GetTagsString(),
			})
		}
	}

	return &templates.TemplType{
		AnvilType:          t.AnvilType,
		ModuleImport:       t.ModuleImport,
		GolangType:         t.GolangType,
		Optional:           t.Optional,
		MapProps:           mapProps,
		ImportsUnorganized: t.GetImportsUnorganized(),
	}
}

func enumToTemplEnum(e *types_parser.Enum) *templates.TemplEnum {
	var values []*templates.TemplEnumValue = nil
	if e.Values != nil {
		biggestName := 0
		for _, v := range e.Values {
			name := fmt.Sprintf("%s_%s", e.GolangName, v.Name)
			lenName := len(name)
			if lenName > biggestName {
				biggestName = lenName
			}
		}

		values = make([]*templates.TemplEnumValue, 0, len(e.Values))

		for _, v := range e.Values {
			values = append(values, &templates.TemplEnumValue{
				Name:    v.Name,
				Spacing: strings.Repeat(" ", biggestName-len(v.Name)),
				Idx:     v.Idx,
				Value:   v.Value,
			})
		}
	}

	var deprecatedValues []*templates.TemplEnumValue = nil
	if e.DeprecatedValues != nil {
		biggestName := 0
		for _, v := range e.DeprecatedValues {
			name := fmt.Sprintf("%s_%s", e.GolangName, v.Name)
			lenName := len(name)
			if lenName > biggestName {
				biggestName = lenName
			}
		}

		deprecatedValues = make([]*templates.TemplEnumValue, 0, len(e.DeprecatedValues))

		for _, v := range e.DeprecatedValues {
			deprecatedValues = append(deprecatedValues, &templates.TemplEnumValue{
				Name:    v.Name,
				Spacing: strings.Repeat(" ", biggestName-len(v.Name)),
				Idx:     v.Idx,
				Value:   v.Value,
			})
		}
	}

	return &templates.TemplEnum{
		GolangName:       e.GolangName,
		GolangType:       e.GolangType,
		Values:           values,
		DeprecatedValues: deprecatedValues,
	}
}
