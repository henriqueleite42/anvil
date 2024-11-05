package schemas

// Generator

type Generator struct {
	Name       string         `yaml:"Name" validate:"required"`
	Version    string         `yaml:"Version" validate:"required"`
	Parameters map[string]any `yaml:"Parameters"`
}

// Config

type Config struct {
	ProjectName  string       `yaml:"ProjectName" validate:"required"`
	AnvilVersion string       `yaml:"AnvilVersion" validate:"required"`
	Schemas      []string     `yaml:"Schemas" validate:"required,min=1"`
	Generators   []*Generator `yaml:"Generators"`
}
