package types_parser

import (
	"fmt"
	"sort"
	"strings"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

type typeParser struct {
	schema *schemas.Schema

	types   map[string]*Type
	enums   map[string]*Enum
	imports map[string]bool
}

func (self *typeParser) ParseType(t *schemas.Type) (*Type, error) {
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

	if existentType, ok := self.types[t.Name]; ok {
		return existentType, nil
	}

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
			GolangType: enum.GolangName,
			AnvilType:  "Enum",
		}
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

		resolvedChildType, err := self.ParseType(childType)
		if err != nil {
			return nil, err
		}

		result = &Type{
			GolangType: "[]" + resolvedChildType.GolangType,
			AnvilType:  "List",
		}
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

		props := make([]*MapProp, len(types), len(types))

		for k, v := range types {
			targetLen := biggest - len(v.Name)

			propType, err := self.ParseType(v)
			if err != nil {
				return nil, err
			}

			resultPropType := propType.GolangType
			if propType.AnvilType == schemas.TypeType_Map {
				resultPropType = "*" + resultPropType
			}

			props[k] = &MapProp{
				Name:    v.Name,
				Spacing: strings.Repeat(" ", targetLen),
				Type:    resultPropType,
			}
		}

		result = &Type{
			GolangType: t.Name,
			AnvilType:  schemas.TypeType_Map,
			MapProps:   props,
		}

		self.types[t.Name] = result
	}

	if t.Optional && t.Type != schemas.TypeType_Map && t.Type != schemas.TypeType_List {
		result.GolangType = "*" + result.GolangType
	}

	return result, nil
}

func (self *typeParser) ParseEnum(e *schemas.Enum) (*Enum, error) {
	if existentEnum, ok := self.enums[e.Name]; ok {
		return existentEnum, nil
	}

	var eType string
	if e.Type == schemas.EnumType_String {
		eType = "string"
	} else {
		eType = "int"
	}

	enum := &Enum{
		GolangName: e.Name,
		GolangType: eType,
		Values:     []*EnumValue{},
	}

	biggest := len(e.Values[0].Name)
	for _, v := range e.Values {
		newLen := len(v.Name)
		if newLen > biggest {
			biggest = newLen
		}
	}

	for k, v := range e.Values {
		targetLen := biggest - len(v.Name)

		enum.Values = append(enum.Values, &EnumValue{
			Idx:     k,
			Name:    v.Name,
			Spacing: strings.Repeat(" ", targetLen),
			Value:   v.Value,
		})
	}

	self.enums[e.Name] = enum

	return enum, nil
}

func (self *typeParser) GetNecessaryImports() []string {
	imports := make([]string, 0, len(self.imports))
	for k := range self.imports {
		imports = append(imports, k)
	}
	return imports
}

func (self *typeParser) GetMapTypes() []*Type {
	types := make([]*Type, 0, len(self.types))
	for _, v := range self.types {
		types = append(types, v)
	}
	sort.Slice(types, func(i, j int) bool {
		return types[i].GolangType < types[j].GolangType
	})
	return types
}

func (self *typeParser) GetEnums() []*Enum {
	enums := make([]*Enum, 0, len(self.enums))
	for _, v := range self.enums {
		enums = append(enums, v)
	}
	sort.Slice(enums, func(i, j int) bool {
		return enums[i].GolangType < enums[j].GolangType
	})
	return enums
}

func NewTypeParser(schema *schemas.Schema) (TypeParser, error) {
	return &typeParser{
		schema: schema,

		types:   map[string]*Type{},
		imports: map[string]bool{},
	}, nil
}
