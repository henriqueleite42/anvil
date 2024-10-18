package templates

import (
	_ "embed"

	"github.com/henriqueleite42/anvil/language-helpers/golang/grpc"
	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

type TemplMethod struct {
	MethodName      string
	MethodNameCamel string
	Input           *grpc.Type
	Output          *grpc.Type
}

type TemplInput struct {
	Domain                      string
	DomainCamel                 string
	DomainSnake                 string
	SpacingRelativeToDomainName string
	ImportsContract             [][]string
	ImportsImplementation       [][]string
	Enums                       []*types_parser.Enum
	Types                       []*types_parser.Type
	Methods                     []*TemplMethod
}

//go:embed contract.txt
var ContractTempl string

//go:embed implementation.txt
var ImplementationTempl string
