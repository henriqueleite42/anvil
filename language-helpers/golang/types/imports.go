package types_parser

import (
	"sort"
	"strings"
)

func (self *typeParser) AddImport(impt string) {
	self.imports[impt] = true
}

func (self *typeParser) GetImports() [][]string {
	// Imports from golang std library
	importsStd := make([]string, 0, len(self.imports))
	// Imports from external libraries
	importsExt := make([]string, 0, len(self.imports))

	for k := range self.imports {
		parts := strings.Split(k, "/")
		if strings.Contains(parts[0], ".") {
			importsExt = append(importsExt, k)
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

func (self *typeParser) ResetImports() {
	self.imports = map[string]bool{}
}
