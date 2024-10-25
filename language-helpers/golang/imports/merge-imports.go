package imports

func (self *importsManager) MergeImport(i *Import) {
	self.AddImport(i.Path, &i.Alias)
}

func (self *importsManager) MergeImports(i []*Import) {
	for _, v := range i {
		self.MergeImport(v)
	}
}
