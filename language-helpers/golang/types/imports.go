package types_parser

import "github.com/henriqueleite42/anvil/language-helpers/golang/imports"

// Get ALL imports from type and child types recursively
func (self *Type) GetImportsUnorganized() []*imports.Import {
	if self.imports == nil {
		return []*imports.Import{}
	}

	return self.imports.GetImportsUnorganized()
}
