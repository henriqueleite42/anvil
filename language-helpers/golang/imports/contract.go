package imports

type Import struct {
	Path            string
	Alias           string
	IsAliasRequired bool
}

// ImportsManager is responsible for helping you to manage Golang Imports
//
// It automatically splits and sorts the imports into standard and external,
// and has support for import aliases
type ImportsManager interface {
	// Add one import to the list
	AddImport(path string, alias *string)
	// Remove one import to the list
	RemoveImport(path string)

	// Create a new import in the current import manager based on an import pointer
	// Doesn't use the same point to avoid any kind of conflict by mutability
	MergeImport(i *Import)
	MergeImports(i []*Import)

	// Return how many imports are in the list
	GetImportsLen() int
	// Return imports without any kind of sorting or filtering
	GetImportsUnorganized() []*Import

	// Return sorted imports, divided into standard and external, and filters the current pkg
	ResolveImports(curPkg string) [][]string
}
