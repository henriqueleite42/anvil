package schemas

// Generator

type Generator struct {
	Name       string         `yaml:"Name"`
	Version    string         `yaml:"Version"`
	Parameters map[string]any `yaml:"Parameters"`
}

// Config

type Config struct {
	AnvilVersion string       `yaml:"AnvilVersion"`
	Schemas      []string     `yaml:"Schemas"`
	Generators   []*Generator `yaml:"Generators"`
}
