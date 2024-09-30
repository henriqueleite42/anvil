package types_parser

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

type typeParser struct {
	schema *schemas.Schema

	typesToAvoidDuplication map[string]*Type
	enumsToAvoidDuplication map[string]*Enum
	types                   []*Type
	enums                   []*Enum
	imports                 map[string]bool
}

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
			// In models file
			golangType = enum.GolangName
		}

		result = &Type{
			GolangType: golangType,
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
		if existentType, ok := self.typesToAvoidDuplication[t.Name]; ok {
			return existentType, nil
		}

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

			propType, err := self.ParseType(v, opt)
			if err != nil {
				return nil, err
			}

			resultPropType := propType.GolangType
			if propType.AnvilType == schemas.TypeType_Map {
				resultPropType = "*" + resultPropType
			}

			prop := &MapProp{
				Name:       v.Name,
				Spacing1:   strings.Repeat(" ", targetLen),
				GolangType: resultPropType,
			}

			if !v.Optional && !slices.Contains(v.Validate, "required") {
				if prop.Tags == nil {
					prop.Tags = []string{}
				}

				prop.Tags = append(prop.Tags, "required")
			}

			if len(v.Validate) > 0 {
				if prop.Tags == nil {
					prop.Tags = []string{}
				}

				prop.Tags = append(prop.Tags, fmt.Sprintf("validate:\"%s\"", strings.Join(v.Validate, ",")))
			}

			props[k] = prop
		}

		if len(props) > 0 {
			biggestType := 0
			for _, v := range props {
				newLen := len(v.GolangType)
				if newLen > biggestType {
					biggestType = newLen
				}
			}

			for _, v := range props {
				targetLen := biggestType - len(v.GolangType)
				v.Spacing2 = strings.Repeat(" ", targetLen)
			}
		}

		result = &Type{
			GolangType: t.Name,
			AnvilType:  schemas.TypeType_Map,
			MapProps:   props,
		}

		self.typesToAvoidDuplication[t.Name] = result
		self.types = append(self.types, result)
	}

	if t.Optional && t.Type != schemas.TypeType_Map && t.Type != schemas.TypeType_List {
		result.GolangType = "*" + result.GolangType
	}

	return result, nil
}

func (self *typeParser) ParseEnum(e *schemas.Enum) (*Enum, error) {
	if existentEnum, ok := self.enumsToAvoidDuplication[e.Name]; ok {
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

	self.enumsToAvoidDuplication[e.Name] = enum
	self.enums = append(self.enums, enum)

	return enum, nil
}

func (self *typeParser) AddImport(impt string) {
	self.imports[impt] = true
}

func (self *typeParser) GetImports() [][]string {
	// Imports from golang std library
	importsStd := make([]string, 0, len(self.imports))
	// Imports from external libraries
	importsExt := make([]string, 0, len(self.imports))

	for k := range self.imports {
		parts := strings.Split(k, "/")
		if strings.Contains(parts[0], ".") {
			importsExt = append(importsExt, k)
		} else {
			importsStd = append(importsStd, k)
		}
	}
	sort.Slice(importsStd, func(i, j int) bool {
		return importsStd[i] < importsStd[j]
	})
	sort.Slice(importsExt, func(i, j int) bool {
		return importsExt[i] < importsExt[j]
	})

	importsResolved := make([][]string, 0, 2)

	if len(importsStd) > 0 {
		importsResolved = append(importsResolved, importsStd)
	}
	if len(importsExt) > 0 {
		importsResolved = append(importsResolved, importsExt)
	}

	return importsResolved
}

func (self *typeParser) ResetImports() {
	self.imports = map[string]bool{}
}

func (self *typeParser) GetMapTypes() []*Type {
	return self.types
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

		typesToAvoidDuplication: map[string]*Type{},
		enumsToAvoidDuplication: map[string]*Enum{},
		types:                   []*Type{},
		enums:                   []*Enum{},
		imports:                 map[string]bool{},
	}, nil
}
