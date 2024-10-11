package types_parser

import (
	"strings"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

type MapProp struct {
	Name     string
	Spacing1 string // Spacing between name and type
	Type     *Type
	Spacing2 string // Spacing between type and tags
	Tags     []string
}

func (self *MapProp) GetTagsString() string {
	return strings.Join(self.Tags, " ")
}

type Type struct {
	GolangPkg  *string // Only Maps, Enums and Lists (because their children can be Maps or Enums) have a pkg
	GolangType string
	AnvilType  schemas.TypeType
	Optional   bool
	MapProps   []*MapProp
}

func (self *Type) GetTypeName(curPkg string) string {
	typeName := self.GolangType

	if self.GolangPkg != nil && *self.GolangPkg != curPkg {
		if self.AnvilType == schemas.TypeType_List {
			trueType := strings.TrimPrefix(self.GolangType, "[]")
			if strings.HasPrefix(trueType, "*") {
				trueType = strings.TrimPrefix(trueType, "*")
				typeName = "[]*" + *self.GolangPkg + "." + trueType
			} else {
				typeName = "[]" + *self.GolangPkg + "." + trueType
			}
		} else {
			typeName = *self.GolangPkg + "." + typeName
		}
	}

	return typeName
}

func (self *Type) GetFullTypeName(curPkg string) string {
	typeName := self.GetTypeName(curPkg)

	if self.AnvilType == schemas.TypeType_Map {
		typeName = "*" + typeName
	}

	if self.Optional &&
		self.AnvilType != schemas.TypeType_Map &&
		self.AnvilType != schemas.TypeType_List {
		typeName = "*" + typeName
	}

	return typeName
}

type EnumValue struct {
	Idx     int32
	Name    string
	Spacing string
	Value   string
}

type Enum struct {
	GolangPkg        string
	GolangName       string
	GolangType       string
	Values           []*EnumValue
	DeprecatedValues []*EnumValue
}

func (self *Enum) GetFullEnumName(curPkg string) string {
	enumName := self.GolangName

	if self.GolangPkg != curPkg {
		enumName = self.GolangPkg + "." + enumName
	}

	return enumName
}

type TypeParser interface {
	// Parse a type and all it's children (if any), then adds them all to the list and returns the root parsed type
	ParseType(t *schemas.Type) (*Type, error)
	// Parse an enum, then adds it to the list and returns the parsed enum
	ParseEnum(e *schemas.Enum) (*Enum, error)

	// Add an import to the list (already handles duplicated imports)
	AddImport(impt string)
	// Returns imports divided by groups (like the formatter does), each group is sorted alphabetically
	GetImports() [][]string
	// Remove all imports from list
	ResetImports()

	// Returns all parsed enums, sorted alphabetically
	GetEnums() []*Enum
	// Returns all parsed types, sorted alphabetically
	GetTypes() []*Type
	// Returns all parsed events, sorted alphabetically
	GetEvents() []*Type
	// Returns all parsed entities, sorted alphabetically
	GetEntities() []*Type
	// Returns all parsed repository types, sorted alphabetically
	GetRepository() []*Type
	// Returns all parsed usecase types, sorted alphabetically
	GetUsecase() []*Type
}
