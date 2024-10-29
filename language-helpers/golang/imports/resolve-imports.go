package imports

import (
	"fmt"
	"sort"
	"strings"
)

func IsExtImport(path string) bool {
	return strings.Contains(path, ".")
}

// - Filters out the current pkg
// - Splits the imports into standard and external
// - Sort both groups alphabetically
// - Resolve import to it's final value to be put into the code, Ex: `alias "foo/bar"`
func ResolveImports(imports []*Import, curModuleAlias string) [][]string {
	// Imports from golang std library
	importsStd := make([]*Import, 0, len(imports))
	// Imports from external libraries
	importsExt := make([]*Import, 0, len(imports))

	for _, v := range imports {
		if v.Alias == curModuleAlias {
			continue
		}

		if IsExtImport(v.Path) {
			importsExt = append(importsExt, v)
		} else {
			importsStd = append(importsStd, v)
		}
	}
	sort.Slice(importsStd, func(i, j int) bool {
		return importsStd[i].Path < importsStd[j].Path
	})
	sort.Slice(importsExt, func(i, j int) bool {
		return importsExt[i].Path < importsExt[j].Path
	})

	// Imports from golang std library
	importsStdString := make([]string, 0, len(importsStd))
	// Imports from external libraries
	importsExtString := make([]string, 0, len(importsExt))

	for _, v := range importsStd {
		if v.IsAliasRequired {
			importsStdString = append(importsStdString, fmt.Sprintf("%s \"%s\"", v.Alias, v.Path))
		} else {
			importsStdString = append(importsStdString, fmt.Sprintf("\"%s\"", v.Path))
		}
	}
	for _, v := range importsExt {
		if v.IsAliasRequired {
			importsExtString = append(importsExtString, fmt.Sprintf("%s \"%s\"", v.Alias, v.Path))
		} else {
			importsExtString = append(importsExtString, fmt.Sprintf("\"%s\"", v.Path))
		}
	}

	importsResolved := make([][]string, 0, 2)

	if len(importsStd) > 0 {
		importsResolved = append(importsResolved, importsStdString)
	}
	if len(importsExt) > 0 {
		importsResolved = append(importsResolved, importsExtString)
	}

	return importsResolved
}
