package schema

type FieldConfidentiality = string

const (
	FieldConfidentiality_Low    FieldConfidentiality = "LOW"
	FieldConfidentiality_Medium FieldConfidentiality = "MEDIUM"
	FieldConfidentiality_High   FieldConfidentiality = "HIGH"
)

type Dependency_Import struct {
	Alias *string `yaml:"Alias,omitempty"`
	Path  string  `yaml:"Path"`
}

type Dependency struct {
	Import *Dependency_Import `yaml:"Import,omitempty"`
	Type   string             `yaml:"Type"`
}

type FieldType = string

const (
	FieldType_String        FieldType = "String"
	FieldType_Int           FieldType = "Int"
	FieldType_Timestamp     FieldType = "Timestamp"
	FieldType_Enum          FieldType = "Enum"
	FieldType_Map           FieldType = "Map"
	FieldType_MapStringMap  FieldType = "Map[String]Map"
	FieldType_ListString    FieldType = "List[String]"
	FieldType_ListInt       FieldType = "List[Int]"
	FieldType_ListTimestamp FieldType = "List[Timestamp]"
	FieldType_ListEnum      FieldType = "List[Enum]"
	FieldType_ListMap       FieldType = "List[Map]"
)

type Field struct {
	Type            FieldType            `yaml:"Type"`
	Confidentiality FieldConfidentiality `yaml:"Confidentiality"`
	Optional        bool                 `yaml:"Optional"`
	Format          *string              `yaml:"Format,omitempty"`
	DbType          *string              `yaml:"DbType,omitempty"`
	Validate        []string             `yaml:"Validate,omitempty"`
	// Used for Map and List[Map]
	Properties map[string]*Field `yaml:"Properties,omitempty"`
	// Used for Enum and List[Enum]
	Values map[string]string `yaml:"Values,omitempty"`
}
