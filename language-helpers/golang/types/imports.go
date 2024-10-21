package types_parser

import (
	"sort"
	"strings"
)

func (self *typeParser) getImports(imports map[string]bool, curPkg string) [][]string {
	// Imports from golang std library
	importsStd := make([]string, 0, len(imports))
	// Imports from external libraries
	importsExt := make([]string, 0, len(imports))

	curPkgImport := self.moduleName + "/" + curPkg

	for k := range imports {
		if strings.Contains(k, ".") {
			if k != curPkgImport {
				importsExt = append(importsExt, k)
			}
		} else {
			importsStd = append(importsStd, k)
		}
	}
	sort.Slice(importsStd, func(i, j int) bool {
		return importsStd[i] < importsStd[j]
	})
	sort.Slice(importsExt, func(i, j int) bool {
		return importsExt[i] < importsExt[j]
	})

	importsResolved := make([][]string, 0, 2)

	if len(importsStd) > 0 {
		importsResolved = append(importsResolved, importsStd)
	}
	if len(importsExt) > 0 {
		importsResolved = append(importsResolved, importsExt)
	}

	return importsResolved
}

func (self *typeParser) AddTypesImport(impt string) {
	self.importsTypes[impt] = true
}
func (self *typeParser) AddEventsImport(impt string) {
	self.importsEvents[impt] = true
}
func (self *typeParser) AddEntitiesImport(impt string) {
	self.importsEntities[impt] = true
}
func (self *typeParser) AddRepositoryImport(impt string) {
	self.importsRepository[impt] = true
}
func (self *typeParser) AddUsecaseImport(impt string) {
	self.importsUsecase[impt] = true
}

func (self *typeParser) GetTypesImports(curPkg string) [][]string {
	return self.getImports(self.importsTypes, curPkg)
}
func (self *typeParser) GetEventsImports(curPkg string) [][]string {
	return self.getImports(self.importsEvents, curPkg)
}
func (self *typeParser) GetEntitiesImports(curPkg string) [][]string {
	return self.getImports(self.importsEntities, curPkg)
}
func (self *typeParser) GetRepositoryImports(curPkg string) [][]string {
	return self.getImports(self.importsRepository, curPkg)
}
func (self *typeParser) GetUsecaseImports(curPkg string) [][]string {
	return self.getImports(self.importsUsecase, curPkg)
}
