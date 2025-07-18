package templates

import (
	_ "embed"

	"github.com/henriqueleite42/anvil/language-helpers/golang/grpc"
	"github.com/henriqueleite42/anvil/language-helpers/golang/imports"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

type TemplEnumValue struct {
	Idx     uint
	Name    string
	Spacing string
	Value   string
}

type TemplEnum struct {
	AnvilEnum *schemas.Enum

	GolangName       string // Enum name
	GolangType       string // string, int, etc
	Values           []*TemplEnumValue
	DeprecatedValues []*TemplEnumValue
}

type UsecaseMethodTemplInput struct {
	DomainPascal   string
	DomainCamel    string
	DomainSnake    string
	MethodName     string
	InputTypeName  string
	OutputTypeName string
	Imports        [][]string
}

type RepositoryMethodTemplInput struct {
	DomainPascal   string
	DomainCamel    string
	DomainSnake    string
	MethodName     string
	InputTypeName  string
	OutputTypeName string
	Imports        [][]string
}

type TemplGrpcMethodDelivery struct {
	DomainPascal string
	DomainCamel  string
	DomainSnake  string
	MethodName   string
	Order        uint
	Input        *grpc.ConvertedValue
	Output       *grpc.ConvertedValue
}

type TemplHttpMethodDelivery struct {
	DomainPascal    string
	DomainCamel     string
	DomainSnake     string
	RouteNamePascal string
	RouteNameSnake  string
	Order           uint
}

type TemplQueueMethodDelivery struct {
	DomainPascal  string
	DomainCamel   string
	MethodName    string
	QueueIdPascal string
	InputName     string
	Input         *types_parser.Type
}

type TemplMethod struct {
	MethodName     string
	InputTypeName  string
	OutputTypeName string
	Order          uint
	Imports        [][]string
}

type GoConfig struct {
	PkgName   string
	GoVersion string
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

type TemplInput struct {
	ProjectName string

	DomainPascal                string
	DomainCamel                 string
	DomainSnake                 string
	SpacingRelativeToDomainName string

	ImportsModels             [][]string
	ImportsRepository         [][]string
	ImportsUsecase            [][]string
	ImportsGrpcDelivery       [][]string
	ImportsGrpcDeliveryHelper [][]string
	ImportsHttpDelivery       [][]string
	ImportsQueueDelivery      [][]string

	Enums    []*TemplEnum
	Types    []*TemplType
	Events   []*TemplType
	Entities []*TemplType

	TypesRepository []*TemplType
	TypesUsecase    []*TemplType

	MethodsRepository    []*TemplMethod
	MethodsUsecase       []*TemplMethod
	MethodsGrpcDelivery  []*TemplGrpcMethodDelivery
	MethodsHttpDelivery  []*TemplHttpMethodDelivery
	MethodsQueueDelivery []*TemplQueueMethodDelivery
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

//go:embed grpc-delivery-module-helper.txt
var GrpcDeliveryModuleHelperTempl string

//go:embed grpc-delivery-module.txt
var GrpcDeliveryModuleTempl string

//go:embed grpc-delivery.txt
var GrpcDeliveryTempl string

//go:embed http-delivery-module.txt
var HttpDeliveryModuleTempl string

//go:embed http-delivery-route.txt
var HttpDeliveryRouteTempl string

//go:embed http-delivery.txt
var HttpDeliveryTempl string

//go:embed main.txt
var MainTempl string

//go:embed models.txt
var ModelsTempl string

//go:embed queue-delivery-module.txt
var QueueDeliveryModuleTempl string

//go:embed queue-delivery.txt
var QueueDeliveryTempl string

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
