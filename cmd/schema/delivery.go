package schema

type Delivery struct {
	Dependencies map[string]*Dependency `yaml:"Dependencies,omitempty"`
}
