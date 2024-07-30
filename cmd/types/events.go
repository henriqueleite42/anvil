package types

type Event struct {
	Formats []string
	Fields  map[string]*Field
}

type Events map[string]*Event
