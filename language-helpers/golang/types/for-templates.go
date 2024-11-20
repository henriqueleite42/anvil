package types_parser

import (
	"strings"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

// ---------------------------------
//
// Utils
//
// ---------------------------------

// Recursively removes [] and * from the prefix of the type,
// so it can handle arrays of arrays
func splitTypePrefix(input string) (string, string) {
	var prefixBuilder strings.Builder
	i := 0

	// Iterate over the string and extract the prefix
	for i < len(input) {
		if strings.HasPrefix(input[i:], "[]") {
			prefixBuilder.WriteString("[]")
			i += 2
		} else if input[i] == '*' {
			prefixBuilder.WriteByte('*')
			i++
		} else {
			break
		}
	}

	// Prefix contains only "[]" and "*", remainder contains the rest
	prefix := prefixBuilder.String()
	remainder := input[i:]

	return prefix, remainder
}

// ---------------------------------
//
// Functions to be used in TEMPLATES
//
// ---------------------------------

func (self *MapProp) GetTagsString() string {
	if self.Tags == nil {
		return ""
	}

	return strings.Join(self.Tags, " ")
}

func (self *Type) GetTypeName(curModuleAlias string) string {
	if self.ModuleImport == nil || self.ModuleImport.Alias == curModuleAlias {
		return self.GolangType
	}

	moduleAlias := self.ModuleImport.Alias

	if self.AnvilType.Type != schemas.TypeType_List {
		return moduleAlias + "." + self.GolangType
	}

	prefix, remainder := splitTypePrefix(self.GolangType)

	return prefix + moduleAlias + "." + remainder
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
