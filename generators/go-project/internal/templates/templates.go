package templates

import (
	_ "embed"

	"github.com/henriqueleite42/anvil/language-helpers/golang/grpc"
	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

type TemplEnumValue struct {
	Idx     int
	Name    string
	Spacing string
	Value   string
}

type TemplEnum struct {
	Name   string
	Type   string
	Values []*TemplEnumValue
}

type UsecaseMethodTemplInput struct {
	Domain         string
	DomainSnake    string
	MethodName     string
	InputTypeName  string
	OutputTypeName string
}

type RepositoryMethodTemplInput struct {
	Domain         string
	DomainSnake    string
	MethodName     string
	InputTypeName  string
	OutputTypeName string
}

type TemplMethodDelivery struct {
	Domain      string
	DomainCamel string
	DomainSnake string
	MethodName  string
	Order       uint
	Input       *grpc.Type
	Output      *grpc.Type
}

type TemplMethod struct {
	MethodName     string
	InputTypeName  string
	OutputTypeName string
	Order          uint
}

type GoConfig struct {
	PkgName   string
	GoVersion string
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
	Events   []*types_parser.Type

	TypesRepository []*types_parser.Type
	TypesUsecase    []*types_parser.Type

	MethodsRepository   []*TemplMethod
	MethodsUsecase      []*TemplMethod
	MethodsGrpcDelivery []*TemplMethodDelivery
}

//go:embed delivery-gratefully-shutdown.txt
var GratefullyShutdownTempl string

//go:embed delivery.txt
var DeliveryTempl string

//go:embed editorconfig.txt
var EditorConfigTempl string

//go:embed gitignore.txt
var GitIgnoreTempl string

//go:embed go-mod.txt
var GoModTempl string

//go:embed grpc-delivery-module.txt
var GrpcDeliveryModuleTempl string

//go:embed grpc-delivery.txt
var GrpcDeliveryTempl string

//go:embed main.txt
var MainTempl string

//go:embed models.txt
var ModelsTempl string

//go:embed readme.txt
var ReadMeTempl string

//go:embed repository-method.txt
var RepositoryMethodTempl string

//go:embed repository-struct.txt
var RepositoryStructTempl string

//go:embed repository.txt
var RepositoryTempl string

//go:embed staticcheck.txt
var StaticCheckTempl string

//go:embed usecase-method.txt
var UsecaseMethodTempl string

//go:embed usecase-struct.txt
var UsecaseStructTempl string

//go:embed usecase.txt
var UsecaseTempl string

//go:embed utils-sync.txt
var UtilsSyncTempl string

//go:embed validator-contract.txt
var ValidatorTempl string

//go:embed validator-implementation.txt
var ValidatorImplementationTempl string
