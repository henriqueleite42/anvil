package schema

type TextCase = string

const (
	TextCase_SnakeCase TextCase = "snake"
)

type Entity_ForeignKey struct {
	Columns    []string `yaml:"Columns"`
	RefTable   string   `yaml:"RefTable"`
	RefColumns []string `yaml:"RefColumn"`
	OnDelete   *string  `yaml:"OnDelete,omitempty"`
	OnUpdate   *string  `yaml:"OnUpdate,omitempty"`
}

type Entity_Index struct {
	Columns []string `yaml:"Columns"`
	Unique  bool     `yaml:"Unique"`
}

type Entity struct {
	Name        *string              `yaml:"Name,omitempty"`
	Columns     map[string]*Field    `yaml:"Columns,omitempty"`
	PrimaryKeys []string             `yaml:"PrimaryKeys,omitempty"`
	Indexes     []*Entity_Index      `yaml:"Indexes,omitempty"`
	ForeignKeys []*Entity_ForeignKey `yaml:"ForeignKeys,omitempty"`
}

type Entities struct {
	Schema      *string            `yaml:"Schema,omitempty"`
	ColumnsCase *string            `yaml:"ColumnsCase,omitempty"`
	Tables      map[string]*Entity `yaml:"Tables,omitempty"`
}
