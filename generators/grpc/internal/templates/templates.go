package templates

import (
	_ "embed"
	"fmt"
)

type ProtofileTemplInputEnumValue struct {
	Name    string
	Spacing string
	Idx     uint
}

type ProtofileTemplInputEnum struct {
	Name             string
	Values           []*ProtofileTemplInputEnumValue
	DeprecatedValues []uint
}

func (self *ProtofileTemplInputEnum) GetDeprecatedValues() string {
	var result string
	for _, v := range self.DeprecatedValues {
		if len(result) > 0 {
			result += fmt.Sprintf(", %v", v)
		} else {
			result += fmt.Sprintf("%v", v)
		}
	}

	return result
}

type ProtofileTemplInputMethod struct {
	Name   string
	Input  *string
	Output *string
}

type ProtofileTemplInputTypeProp struct {
	Type     string
	Spacing1 string // Spacing between type and name
	Name     string
	Spacing2 string // Spacing between name and idx
	Idx      int
}

type ProtofileTemplInputType struct {
	Name  string
	Props []*ProtofileTemplInputTypeProp
}

type ProtofileTemplInput struct {
	Domain  string
	Imports []string
	Methods []*ProtofileTemplInputMethod
	Enums   []*ProtofileTemplInputEnum
	Types   []*ProtofileTemplInputType
}

//go:embed protofile.txt
var ProtofileTempl string
