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
		}
	}
	if t.Type == schemas.TypeType_Int {
		result = &Type{
			GolangType: "int32",
		}
	}
	if t.Type == schemas.TypeType_Float {
		result = &Type{
			GolangType: "float32",
		}
	}
	if t.Type == schemas.TypeType_Bool {
		result = &Type{
			GolangType: "bool",
		}
	}
	if t.Type == schemas.TypeType_Timestamp {
		self.imports["time"] = true
		result = &Type{
			GolangType: "time.Time",
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

		result = &Type{
			GolangPkg:  &enum.GolangPkg,
			GolangType: enum.GolangName,
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

			var childOpt *ParseTypeOpt = opt
			if opt == nil {
				childOpt = &ParseTypeOpt{}
			}
			childOpt.prefixForChildren = t.Name

			propType, err := self.ParseType(childType, childOpt)
			if err != nil {
				return nil, err
			}

			prop := &MapProp{
				Name: *v.PropName,
				Type: propType,
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

				newLenType := len(v.Type.GolangType)
				if newLenType > biggestType {
					biggestType = newLenType
				}
			}

			for _, v := range props {
				targetLenName := biggestName - len(v.Name)
				v.Spacing1 = strings.Repeat(" ", targetLenName)

				targetLenType := biggestType - len(v.Type.GolangType)
				v.Spacing2 = strings.Repeat(" ", targetLenType)
			}
		}

		golangType := t.Name
		if opt != nil {
			golangType = opt.prefixForChildren + golangType
		}

		result = &Type{
			GolangType: golangType,
			MapProps:   props,
		}

		if strings.HasPrefix(t.Ref, "Types") {
			result.GolangPkg = &self.typesPkg
			self.types = append(self.types, result)
		} else if strings.HasPrefix(t.Ref, "Events") {
			result.GolangPkg = &self.eventsPkg
			self.events = append(self.events, result)
		} else if strings.HasPrefix(t.Ref, "Entities") {
			result.GolangPkg = &self.entitiesPkg
			self.entities = append(self.entities, result)
		} else if strings.HasPrefix(t.Ref, "Repository") {
			result.GolangPkg = &self.repositoryPkg
			self.repository = append(self.repository, result)
		} else if strings.HasPrefix(t.Ref, "Usecase") {
			result.GolangPkg = &self.usecasePkg
			self.usecase = append(self.usecase, result)
		} else {
			return nil, fmt.Errorf("unable to get package for \"%s\"", t.Ref)
		}

		self.typesToAvoidDuplication[t.Ref] = result
	}

	result.AnvilType = t.Type
	result.Optional = t.Optional

	return result, nil
}
