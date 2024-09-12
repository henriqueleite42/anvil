package templates

import "github.com/henriqueleite42/anvil/cli/schemas"

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

type TemplTypeProp struct {
	Name    string
	Spacing string
	Type    string
}

type TemplType struct {
	Name         string
	OriginalType schemas.TypeType
	Props        []*TemplTypeProp
}

type TemplMethodProp struct {
	Name    string
	Spacing string
	Value   string
}

type TemplMethod struct {
	MethodName         string
	MethodNameCamel    string
	Input              []*TemplMethodProp
	InputPropsPrepare  []string
	Output             []*TemplMethodProp
	OutputPropsPrepare []string
}

type TemplInput struct {
	Domain                      string
	DomainCamel                 string
	DomainSnake                 string
	SpacingRelativeToDomainName string
	ImportsContract             []string
	ImportsImplementation       []string
	Enums                       []*TemplEnum
	Types                       []*TemplType
	Methods                     []*TemplMethod
}
