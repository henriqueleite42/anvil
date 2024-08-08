package schema

type Relationship struct {
	Uri string `yaml:"Uri"`
}

type Relationships map[string]*Relationship
