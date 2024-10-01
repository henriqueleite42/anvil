package grpc

import (
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

type Prop struct {
	Name    string
	Spacing string
	Value   string
}

type Type struct {
	Name         string
	Props        []*Prop
	PropsPrepare []string
}

type GoToProtoInput struct {
	Type                    *schemas.Type
	MethodName              string
	VariableName            string
	PrefixForVariableNaming string
	HasOutput               bool
}

type ProtoToGoInput struct {
	Type                    *schemas.Type
	MethodName              string
	VariableName            string
	PrefixForVariableNaming string
	PkgForEnums             string
	HasOutput               bool
}

type GrpcParser interface {
	GoToProto(t *GoToProtoInput) (*Type, error)
	ProtoToGo(t *ProtoToGoInput) (*Type, error)
}
