package types

type FieldConfidentiality = string

const (
	FieldConfidentiality_Low    FieldConfidentiality = "LOW"
	FieldConfidentiality_Medium FieldConfidentiality = "MEDIUM"
	FieldConfidentiality_High   FieldConfidentiality = "HIGH"
)

type Dependency_Import struct {
	Alias *string
	Path  string
}

type Dependency struct {
	Import  *Dependency_Import
	Type    string
	Private bool
}

type FieldType = string

const (
	FieldType_String     FieldType = "String"
	FieldType_Int        FieldType = "Int"
	FieldType_Enum       FieldType = "Enum"
	FieldType_ListString FieldType = "List[String]"
	FieldType_ListInt    FieldType = "List[Int]"
	FieldType_ListEnum   FieldType = "List[Enum]"
	FieldType_List       FieldType = "List"
	FieldType_Map        FieldType = "Map"
)

type Field struct {
	Type            FieldType
	Confidentiality FieldConfidentiality
	Optional        bool
	DbType          *string
	Validate        []string
	// Used for Map
	Properties map[string]*Field
	// Used for List
	Items map[string]*Field
	// Used for Enum and List[Enum]
	Values map[string]string
}
