package types_parser

import "github.com/henriqueleite42/anvil/language-helpers/golang/schemas"

type MapProp struct {
	Name       string
	Spacing1   string // Spacing between name and type
	GolangType string
	Spacing2   string // Spacing between type and tags
	Tags       []string
}

type Type struct {
	GolangType string
	AnvilType  schemas.TypeType
	MapProps   []*MapProp
}

type EnumValue struct {
	Idx     int
	Name    string
	Spacing string
	Value   string
}

type Enum struct {
	GolangName string
	GolangType string
	Values     []*EnumValue
}

type ParseTypeOpt struct {
	PrefixForEnums string
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

	// Returns all parsed map types, sorted by parse order
	GetMapTypes() []*Type
	// Returns all parsed enums, sorted alphabetically
	GetEnums() []*Enum
}
