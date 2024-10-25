package grpc

import (
	"github.com/henriqueleite42/anvil/language-helpers/golang/imports"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

// Input for internal use, while it still is being processes
type convertingInput struct {
	input *ConverterInput

	// Internal use
	overwriteTypeName       *string // For map child types, to add the prefix of their parent
	prefixForVariableNaming *string // To avoid naming duplication
	indentationLvl          int     // Recursive indentation lvl, to keep things formatted
}

// For internal use, value while it still being converted
type convertingValue struct {
	GolangType     string   // Already includes package, and * if necessary. Ex: *foo.Bar, *string, int32
	GolangTypeName string   // Only includes the type name and package. Ex: foo.Bar, string, int32
	ProtoType      string   // Already includes package, and * if necessary. Ex: *foo.Bar, *string, int32
	ProtoTypeName  string   // Only includes the type name and package. Ex: foo.Bar, string, int32
	Value          string   // The value to be used
	Prepare        []string // Things necessary to prepare the values

	imports imports.ImportsManager
}

type ConverterInput struct {
	CurModuleImport *imports.Import // Cur module import representation
	PbModuleImport  *imports.Import // Protobuff module import representation
	Type            *schemas.Type   // Type to convert
	VarToConvert    string          // Variable name to access the value to be converted
}

// Value already converted
type ConvertedValue struct {
	Value              string   // The value to be used. It can be the name of an variable or the value directly
	Prepare            []string // Things necessary to prepare the value. In the template, it must come BEFORE the usage of the value
	ImportsUnorganized []*imports.Import
}

type GrpcParser interface {
	// Doesn't save any kind of state, can be used multiple times for different
	// types and they will not affect each other
	GoToProto(i *ConverterInput) (*ConvertedValue, error)
	// Doesn't save any kind of state, can be used multiple times for different
	// types and they will not affect each other
	ProtoToGo(i *ConverterInput) (*ConvertedValue, error)
}
