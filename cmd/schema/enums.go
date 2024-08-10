package schema

type EnumValue struct {
	Name  string `yaml:"Name"`
	Value string `yaml:"Value"`
}

type Enum struct {
	Type   string       `yaml:"Type"`
	Values []*EnumValue `yaml:"Values"`
}

type Enums map[string]*Enum
