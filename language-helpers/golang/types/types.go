package types_parser

import (
	"strings"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

type MapProp struct {
	Name     string
	Spacing1 string // Spacing between name and type
	Type     *Type
	Spacing2 string // Spacing between type and tags
	Tags     []string
}

func (self *MapProp) GetTagsString() string {
	return strings.Join(self.Tags, " ")
}

type Type struct {
	GolangPkg  *string // Only Maps, Enums and Lists (because their children can be Maps or Enums) have a pkg
	GolangType string
	AnvilType  schemas.TypeType
	Optional   bool
	MapProps   []*MapProp
}

func (self *Type) GetTypeName(curPkg string) string {
	typeName := self.GolangType

	if self.GolangPkg != nil && *self.GolangPkg != curPkg {
		if self.AnvilType == schemas.TypeType_List {
			trueType := strings.TrimPrefix(self.GolangType, "[]")
			if strings.HasPrefix(trueType, "*") {
				trueType = strings.TrimPrefix(trueType, "*")
				typeName = "[]*" + *self.GolangPkg + "." + trueType
			} else {
				typeName = "[]" + *self.GolangPkg + "." + trueType
			}
		} else {
			typeName = *self.GolangPkg + "." + typeName
		}
	}

	return typeName
}

func (self *Type) GetFullTypeName(curPkg string) string {
	typeName := self.GetTypeName(curPkg)

	if self.AnvilType == schemas.TypeType_Map {
		typeName = "*" + typeName
	}

	if self.Optional &&
		self.AnvilType != schemas.TypeType_Map &&
		self.AnvilType != schemas.TypeType_List {
		typeName = "*" + typeName
	}

	return typeName
}

type EnumValue struct {
	Idx     uint
	Name    string
	Spacing string
	Value   string
}

type Enum struct {
	GolangPkg        string
	GolangName       string
	GolangType       string
	Values           []*EnumValue
	DeprecatedValues []*EnumValue
}

func (self *Enum) GetFullEnumName(curPkg string) string {
	enumName := self.GolangName

	if self.GolangPkg != curPkg {
		enumName = self.GolangPkg + "." + enumName
	}

	return enumName
}
