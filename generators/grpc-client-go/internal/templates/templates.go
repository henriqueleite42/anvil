package templates

import (
	_ "embed"

	"github.com/henriqueleite42/anvil/language-helpers/golang/grpc"
	"github.com/henriqueleite42/anvil/language-helpers/golang/imports"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

type TemplEnumValue struct {
	Idx     uint
	Name    string
	Spacing string
	Value   string
}

type TemplEnum struct {
	GolangName       string // Enum name
	GolangType       string // string, int, etc
	Values           []*TemplEnumValue
	DeprecatedValues []*TemplEnumValue
}

type TemplTypeMapProp struct {
	Name     string
	Spacing1 string
	Type     string
	Spacing2 string
	Tags     string
}

type TemplType struct {
	AnvilType *schemas.Type

	ModuleImport *imports.Import // Import of the module of the type, only Maps, Enums and Lists (of Maps nad Enums) will have one
	GolangType   string
	Optional     bool
	MapProps     []*TemplTypeMapProp

	ImportsUnorganized []*imports.Import
}

type TemplMethod struct {
	MethodName      string
	MethodNameCamel string
	Input           *grpc.ConvertedValue
	Output          *grpc.ConvertedValue
}

type TemplInput struct {
	DomainPascal                string
	DomainCamel                 string
	DomainSnake                 string
	SpacingRelativeToDomainName string
	ImportsContract             [][]string
	ImportsImplementation       [][]string
	Enums                       []*TemplEnum
	Types                       []*TemplType
	Methods                     []*TemplMethod
}

//go:embed contract.txt
var ContractTempl string

//go:embed implementation.txt
var ImplementationTempl string
