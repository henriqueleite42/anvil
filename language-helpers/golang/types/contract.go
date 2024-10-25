package types_parser

import (
	"github.com/henriqueleite42/anvil/language-helpers/golang/imports"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

type MapProp struct {
	Name string
	Type *Type
	Tags []string
}

type Type struct {
	ModuleImport *imports.Import // Import of the module of the type, only Maps, Enums and Lists (of Maps nad Enums) will have one
	GolangType   string
	AnvilType    schemas.TypeType
	Optional     bool
	MapProps     []*MapProp

	// Internal use, handles all the imports necessary to
	// create the struct declaration of the type
	//
	// Ex: Id it's a Map that has a prop with type `time.Time`, it will add
	// "time" to the imports list.
	//
	// It doesn't work atr a deep level: If it's a Map with a child Map,
	// The imports will not have the necessary imports for the child Map or
	// the import for the child's module. They will be at the child's type.
	imports imports.ImportsManager
}

type EnumValue struct {
	Idx   uint
	Name  string
	Value string
}

type Enum struct {
	Import           *imports.Import
	GolangName       string // Enum name
	GolangType       string // string, int, etc
	Values           []*EnumValue
	DeprecatedValues []*EnumValue
}

// The objective of the types parser is to help you to convert
// Anvil types to Golang types.
//
// It doesn't parse only special types (structs), but ALL types like
// `string`, `int`, etc
//
// A single instance of types parser can and should be used per domain.
//
// Note that the types parser must be used **per domain**, and not per schema,
// since a single schema can have multiple domains.
type TypesParser interface {
	// Parse an enum, then adds it to the list and returns the parsed enum
	ParseEnum(e *schemas.Enum) (*Enum, error)
	// Parse a type and all it's children (if any), then adds them all to the list and returns the root parsed type
	//
	// If you parse the same type twice, it will simple return the same result, but without parse it again
	ParseType(t *schemas.Type) (*Type, error)

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
