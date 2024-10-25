package imports

func NewImport(path string, alias *string) *Import {
	defaultAlias := GetDefaultAlias(path)

	var isAliasRequired bool
	var finalAlias string
	if alias != nil && *alias != defaultAlias {
		finalAlias = *alias
		isAliasRequired = true
	} else {
		finalAlias = defaultAlias
	}

	return &Import{
		Path:            path,
		Alias:           finalAlias,
		IsAliasRequired: isAliasRequired,
	}
}

func (self *importsManager) AddImport(path string, alias *string) {
	if _, ok := self.imports[path]; ok {
		return
	}

	impt := NewImport(path, alias)

	self.imports[path] = impt
}

func (self *importsManager) RemoveImport(path string) {
	delete(self.imports, path)
}
