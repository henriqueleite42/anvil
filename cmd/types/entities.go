package types

type TextCase = string

const (
	TextCase_SnakeCase TextCase = "snake"
)

type Entity_Index struct {
	Columns []string
	Unique  bool
}

type Entity struct {
	Name        *string
	Columns     map[string]*Field
	PrimaryKeys []string
	Indexes     []*Entity_Index
}

type Entities struct {
	Schema      *string
	ColumnsCase *string
	Tables      map[string]*Entity
}
