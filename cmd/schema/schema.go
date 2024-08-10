package schema

type Schema struct {
	Domain        string         `yaml:"Domain"`
	Relationships *Relationships `yaml:"Relationships,omitempty"`
	Types         *Types         `yaml:"Types,omitempty"`
	Enums         *Enums         `yaml:"Enums,omitempty"`
	Events        *Events        `yaml:"Events,omitempty"`
	Entities      *Entities      `yaml:"Entities,omitempty"`
	Repository    *Repository    `yaml:"Repository,omitempty"`
	Usecase       *Usecase       `yaml:"Usecase,omitempty"`
	Delivery      *Delivery      `yaml:"Delivery,omitempty"`
}
