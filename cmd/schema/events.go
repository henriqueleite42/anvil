package schema

type Event struct {
	Formats []string          `yaml:"Formats,omitempty"`
	Fields  map[string]*Field `yaml:"Fields,omitempty"`
}

type Events map[string]*Event
