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
	GolangPkg  *string // Only Maps and Enums have a pkg
	GolangType string
	AnvilType  schemas.TypeType
	Optional   bool
	MapProps   []*MapProp
}

func (self *Type) GetFullTypeName(curPkg string) string {
	typeName := self.GolangType

	if self.GolangPkg != nil && *self.GolangPkg != curPkg {
		typeName = *self.GolangPkg + "." + typeName
	}

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
	Idx     int
	Name    string
	Spacing string
	Value   string
}

type Enum struct {
	GolangPkg  string
	GolangName string
	GolangType string
	Values     []*EnumValue
}

func (self *Enum) GetFullEnumName(curPkg string) string {
	enumName := self.GolangName

	if self.GolangPkg != curPkg {
		enumName = self.GolangPkg + "." + enumName
	}

	return enumName
}

type ParseTypeOpt struct {
	prefixForChildren string // Internal use, for maps of maps
}

type TypeParser interface {
	// Parse a type and all it's children (if any), then adds them all to the list and returns the root parsed type
	ParseType(t *schemas.Type, opt *ParseTypeOpt) (*Type, error)
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
	// Returns all parsed types, sorted by parse order
	GetTypes() []*Type
	// Returns all parsed events, sorted by parse order
	GetEvents() []*Type
	// Returns all parsed entities, sorted by parse order
	GetEntities() []*Type
	// Returns all parsed repository types, sorted by parse order
	GetRepository() []*Type
	// Returns all parsed usecase types, sorted by parse order
	GetUsecase() []*Type
}
