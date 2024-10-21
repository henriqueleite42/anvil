package types_parser

import (
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

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
	// Add an import to the list (already handles duplicated imports)
	AddTypesImport(impt string)
	// Returns imports divided by groups (like the formatter does), each group is sorted alphabetically
	GetTypesImports(curPkg string) [][]string
	// Add an import to the list (already handles duplicated imports)
	AddEventsImport(impt string)
	// Returns imports divided by groups (like the formatter does), each group is sorted alphabetically
	GetEventsImports(curPkg string) [][]string
	// Add an import to the list (already handles duplicated imports)
	AddEntitiesImport(impt string)
	// Returns imports divided by groups (like the formatter does), each group is sorted alphabetically
	GetEntitiesImports(curPkg string) [][]string
	// Add an import to the list (already handles duplicated imports)
	AddRepositoryImport(impt string)
	// Returns imports divided by groups (like the formatter does), each group is sorted alphabetically
	GetRepositoryImports(curPkg string) [][]string
	// Add an import to the list (already handles duplicated imports)
	AddUsecaseImport(impt string)
	// Returns imports divided by groups (like the formatter does), each group is sorted alphabetically
	GetUsecaseImports(curPkg string) [][]string

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
