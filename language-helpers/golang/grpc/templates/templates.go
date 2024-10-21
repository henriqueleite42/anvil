package templates

import (
	_ "embed"
	"strings"
)

type InputPropOptionalTemplInput struct {
	VarName              string
	OriginalVariableName string
	Type                 string
	ValueToAssign        string
	NeedsPointer         bool
}

//go:embed input-prop-optional.txt
var InputPropOptionalTempl string

type InputPropMapTemplProp struct {
	Name    string
	Spacing string
	Value   string
}

type InputPropMapTemplInput struct {
	IndentationLvl int // internal use, controls sub levels

	HasOutput            bool
	MethodName           string
	Optional             bool
	OriginalVariableName string
	Prepare              []string
	Props                []*InputPropMapTemplProp
	Type                 string
	TypePkg              *string
	VarName              string
}

func (self *InputPropMapTemplInput) Idt() string {
	return strings.Repeat("	", self.IndentationLvl)
}

//go:embed input-prop-map.txt
var InputPropMapTempl string

type InputPropListTemplInput struct {
	IndentationLvl int // internal use, controls sub levels

	ChildOptional        bool
	HasOutput            bool
	MethodName           string
	Optional             bool
	OriginalVariableName string
	Prepare              []string
	Type                 string
	ValueToAppend        string
	VarName              string
}

func (self *InputPropListTemplInput) Idt() string {
	return strings.Repeat("	", self.IndentationLvl)
}

//go:embed input-prop-list.txt
var InputPropListTempl string
