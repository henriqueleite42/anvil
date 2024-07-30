package types

type RepositoryMethod struct {
	Input  map[string]*Field
	Output map[string]*Field
}

type Repository struct {
	Dependencies map[string]*Dependency
	Inputs       map[string]*Dependency
	Methods      map[string]*RepositoryMethod
}
