package grpc

import (
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

type Type struct {
	GolangType     string   // Already includes package, and * if necessary. Ex: *foo.Bar, *string, int32
	GolangTypeName string   // Only includes the type name and package. Ex: foo.Bar, string, int32
	ProtoType      string   // Already includes package, and * if necessary. Ex: *foo.Bar, *string, int32
	ProtoTypeName  string   // Only includes the type name and package. Ex: foo.Bar, string, int32
	Value          string   // The value to bbe used
	Prepare        []string // Things necessary to prepare the values
}

type GoToProtoInput struct {
	Type                     *schemas.Type
	TypeName                 string
	MethodName               string
	VariableToAccessTheValue string
	PrefixForVariableNaming  string
	HasOutput                bool
	CurPkg                   string
	CurDomain                string

	indentationLvl int // Internal use, used for child types
}

type ProtoToGoInput struct {
	Type                     *schemas.Type
	TypeName                 string
	MethodName               string
	VariableToAccessTheValue string
	PrefixForVariableNaming  string
	HasOutput                bool
	CurPkg                   string
	CurDomain                string

	indentationLvl int // Internal use, used for child types
}

type GrpcParser interface {
	GetImports() [][]string

	GetProtoTypeName(curDomain string, t *schemas.Type) (string, error)

	GoToProto(i *GoToProtoInput) (*Type, error)
	ProtoToGo(i *ProtoToGoInput) (*Type, error)
}
