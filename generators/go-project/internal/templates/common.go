package templates

import "github.com/henriqueleite42/anvil/cli/schemas"

type TemplTypeProp struct {
	Name     string
	Type     string
	Tags     string
	Spacing1 string
	Spacing2 string
}

type TemplType struct {
	Name         string
	OriginalType schemas.TypeType
	Props        []*TemplTypeProp
}

type TemplMethod struct {
	MethodName     string
	InputTypeName  string
	OutputTypeName string
}

type TemplInput struct {
	Domain            string
	DomainSnake       string
	ImportsModels     []string
	ImportsRepository []string
	ImportsUsecase    []string
	Enums             []*TemplEnum
	Entities          []*TemplType
	TypesRepository   []*TemplType
	TypesUsecase      []*TemplType
	MethodsRepository []*TemplMethod
	MethodsUsecase    []*TemplMethod
}
