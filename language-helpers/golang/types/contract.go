package types_parser

import "github.com/henriqueleite42/anvil/language-helpers/golang/schemas"

type MapProp struct {
	Name    string
	Spacing string // Spacing to write type in file
	Type    string
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

type TypeParser interface {
	ParseType(t *schemas.Type) (*Type, error)
	ParseEnum(e *schemas.Enum) (*Enum, error)
	GetNecessaryImports() []string
	GetMapTypes() []*Type
	GetEnums() []*Enum
}
