package types_parser

import (
	"fmt"
	"slices"
	"strings"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *typeParser) ParseType(t *schemas.Type, opt *ParseTypeOpt) (*Type, error) {
	var result *Type

	// ----------------------
	//
	// Basic types
	//
	// ----------------------

	if t.Type == schemas.TypeType_String {
		result = &Type{
			GolangType: "string",
			AnvilType:  schemas.TypeType_String,
		}
	}
	if t.Type == schemas.TypeType_Int {
		result = &Type{
			GolangType: "int32",
			AnvilType:  schemas.TypeType_Int,
		}
	}
	if t.Type == schemas.TypeType_Float {
		result = &Type{
			GolangType: "float32",
			AnvilType:  schemas.TypeType_Float,
		}
	}
	if t.Type == schemas.TypeType_Bool {
		result = &Type{
			GolangType: "bool",
			AnvilType:  schemas.TypeType_Bool,
		}
	}
	if t.Type == schemas.TypeType_Timestamp {
		self.imports["time"] = true
		result = &Type{
			GolangType: "time.Time",
			AnvilType:  schemas.TypeType_Timestamp,
		}
	}

	// ----------------------
	//
	// Complex types
	//
	// ----------------------

	if t.Type == schemas.TypeType_Enum {
		if t.EnumHash == nil {
			return nil, fmt.Errorf("enum \"%s\" not found", *t.EnumHash)
		}

		schemaEnum := self.schema.Enums.Enums[*t.EnumHash]
		enum, err := self.ParseEnum(schemaEnum)
		if err != nil {
			return nil, err
		}

		var golangType string
		if opt != nil && opt.PrefixForEnums != "" {
			golangType = fmt.Sprintf("%s.%s", opt.PrefixForEnums, enum.GolangName)
		} else {
			golangType = enum.GolangName
		}

		result = &Type{
			GolangType: golangType,
			AnvilType:  "Enum",
		}
	}
	if t.Type == schemas.TypeType_List {
		if t.ChildTypes == nil {
			return nil, fmt.Errorf("ChildTypes for \"%s\" not found", t.Name)
		}
		if len(t.ChildTypes) != 1 {
			return nil, fmt.Errorf("ChildTypes for \"%s\" must have exactly one item", t.Name)
		}

		childType, ok := self.schema.Types.Types[t.ChildTypes[0].TypeHash]
		if !ok {
			return nil, fmt.Errorf("type \"%s\" not found", t.ChildTypes[0].TypeHash)
		}

		resolvedChildType, err := self.ParseType(childType, opt)
		if err != nil {
			return nil, err
		}

		result = &Type{
			GolangType: "[]" + resolvedChildType.GolangType,
			AnvilType:  "List",
		}
	}

	// ----------------------
	//
	// Maps (child types)
	//
	// ----------------------

	if t.Type == schemas.TypeType_Map {
		if existentType, ok := self.typesToAvoidDuplication[t.Ref]; ok {
			return existentType, nil
		}

		props := make([]*MapProp, len(t.ChildTypes), len(t.ChildTypes))

		for k, v := range t.ChildTypes {
			if v.PropName == nil {
				return nil, fmt.Errorf("ChildType \"%s.%d\" must have a PropName", t.Name, k)
			}

			childType, ok := self.schema.Types.Types[v.TypeHash]
			if !ok {
				return nil, fmt.Errorf("type \"%s\" not found", v.TypeHash)
			}

			propType, err := self.ParseType(childType, opt)
			if err != nil {
				return nil, err
			}

			resultPropType := propType.GolangType
			if propType.AnvilType == schemas.TypeType_Map {
				resultPropType = "*" + resultPropType
			}

			prop := &MapProp{
				Name:       *v.PropName,
				GolangType: resultPropType,
			}

			if !childType.Optional && !slices.Contains(childType.Validate, "required") {
				childType.Validate = append(childType.Validate, "required")
			}

			if len(childType.Validate) > 0 {
				if prop.Tags == nil {
					prop.Tags = []string{}
				}

				prop.Tags = append(prop.Tags, fmt.Sprintf("validate:\"%s\"", strings.Join(childType.Validate, ",")))
			}

			if childType.DbName != nil {
				prop.Tags = append(prop.Tags, fmt.Sprintf("db:\"%s\"", *childType.DbName))
			}

			props[k] = prop
		}

		if len(props) > 0 {
			biggestName := 0
			biggestType := 0
			for _, v := range props {
				newLenName := len(v.Name)
				if newLenName > biggestName {
					biggestName = newLenName
				}

				newLenType := len(v.GolangType)
				if newLenType > biggestType {
					biggestType = newLenType
				}
			}

			for _, v := range props {
				targetLenName := biggestName - len(v.Name)
				v.Spacing1 = strings.Repeat(" ", targetLenName)

				targetLenType := biggestType - len(v.GolangType)
				v.Spacing2 = strings.Repeat(" ", targetLenType)
			}
		}

		result = &Type{
			GolangType: t.Name,
			AnvilType:  schemas.TypeType_Map,
			MapProps:   props,
		}

		self.typesToAvoidDuplication[t.Ref] = result
		self.types = append(self.types, result)
	}

	if t.Optional && t.Type != schemas.TypeType_Map && t.Type != schemas.TypeType_List {
		result.GolangType = "*" + result.GolangType
	}

	return result, nil
}
