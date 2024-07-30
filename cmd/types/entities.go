package types

type TextCase = string

const (
	TextCase_SnakeCase TextCase = "snake"
)

type Entity_ForeignKey struct {
	Column    string
	RefTable  string
	RefColumn string
	OnDelete  *string
	OnUpdate  *string
}

type Entity_Index struct {
	Columns []string
	Unique  bool
}

type Entity struct {
	Name        *string
	Columns     map[string]*Field
	PrimaryKeys []string
	Indexes     []*Entity_Index
	ForeignKeys []*Entity_ForeignKey
}

type Entities struct {
	Schema      *string
	ColumnsCase *string
	Tables      map[string]*Entity
}
