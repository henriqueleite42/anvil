package imports

func (self *importsManager) GetImportsLen() int {
	return len(self.imports)
}

func (self *importsManager) GetImportsUnorganized() []*Import {
	allImports := make([]*Import, 0, len(self.imports))

	for _, v := range self.imports {
		allImports = append(allImports, v)
	}

	return allImports
}

func (self *importsManager) ResolveImports(curPkg string) [][]string {
	return ResolveImports(self.GetImportsUnorganized(), curPkg)
}
