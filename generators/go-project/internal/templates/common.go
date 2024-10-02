package templates

import (
	"github.com/henriqueleite42/anvil/language-helpers/golang/grpc"
	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

type TemplMethodDelivery struct {
	Domain      string
	DomainCamel string
	DomainSnake string
	MethodName  string
	Input       *grpc.Type
	Output      *grpc.Type
}

type TemplMethod struct {
	MethodName     string
	InputTypeName  string
	OutputTypeName string
}

type TemplInput struct {
	Domain                      string
	DomainCamel                 string
	DomainSnake                 string
	SpacingRelativeToDomainName string

	ImportsModels       [][]string
	ImportsRepository   [][]string
	ImportsUsecase      [][]string
	ImportsGrpcDelivery [][]string

	Enums    []*types_parser.Enum
	Entities []*types_parser.Type

	TypesRepository []*types_parser.Type
	TypesUsecase    []*types_parser.Type

	MethodsRepository   []*TemplMethod
	MethodsUsecase      []*TemplMethod
	MethodsGrpcDelivery []*TemplMethodDelivery
}
