package types_parser

import (
	"strings"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

// ---------------------------------
//
// Functions to be used in TEMPLATES
//
// ---------------------------------

func (self *MapProp) GetTagsString() string {
	return strings.Join(self.Tags, " ")
}

func (self *Type) GetTypeName(curModuleAlias string) string {
	if self.ModuleImport.Alias == curModuleAlias {
		return self.GolangType
	}

	moduleAlias := self.ModuleImport.Alias

	if self.AnvilType.Type != schemas.TypeType_List {
		return moduleAlias + "." + self.GolangType
	}

	prefix := "[]"
	trueType := strings.TrimPrefix(self.GolangType, "[]")
	if strings.HasPrefix(trueType, "*") {
		trueType = strings.TrimPrefix(trueType, "*")
		prefix = prefix + "*"
	}

	return prefix + moduleAlias + "." + trueType
}

func (self *Type) GetFullTypeName(curModuleAlias string) string {
	typeName := self.GetTypeName(curModuleAlias)

	if self.AnvilType.Type == schemas.TypeType_Map {
		typeName = "*" + typeName
	}

	if self.Optional &&
		self.AnvilType.Type != schemas.TypeType_Map &&
		self.AnvilType.Type != schemas.TypeType_List {
		typeName = "*" + typeName
	}

	return typeName
}

func (self *Enum) GetFullEnumName(curModuleAlias string) string {
	if self.Import.Path == curModuleAlias {
		return self.GolangName
	}

	return self.Import.Alias + "." + self.GolangType
}
