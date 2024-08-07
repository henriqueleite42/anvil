package schema

type RepositoryMethod struct {
	Input  map[string]*Field `yaml:"Input,omitempty"`
	Output map[string]*Field `yaml:"Output,omitempty"`
}

type Repository struct {
	Dependencies map[string]*Dependency       `yaml:"Dependencies,omitempty"`
	Inputs       map[string]*Dependency       `yaml:"Inputs,omitempty"`
	Methods      map[string]*RepositoryMethod `yaml:"Methods,omitempty"`
}
