package schema

// Common

type FieldType string

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

type FieldConfidentiality string

const (
	FieldConfidentiality_Low    FieldConfidentiality = "LOW"
	FieldConfidentiality_Medium FieldConfidentiality = "MEDIUM"
	FieldConfidentiality_High   FieldConfidentiality = "HIGH"
)

type Field struct {
	Name            string               `yaml:"Name"`
	RootNode        string               `yaml:"RootNode"`
	OriginalPath    string               `yaml:"OriginalPath"`
	StateHash       string               `yaml:"StateHash"`
	Confidentiality FieldConfidentiality `yaml:"Confidentiality"`
	Optional        bool                 `yaml:"Optional"`
	Format          *string              `yaml:"Format,omitempty"`
	Validate        []string             `yaml:"Validate,omitempty"`
	Type            FieldType            `yaml:"Type"`
	DbType          *string              `yaml:"DbType,omitempty"`
	// Used for Map and List[Map]
	Properties map[string]*Field `yaml:"Properties,omitempty"`
	// Used for Enum and List[Enum]
	Values map[string]string `yaml:"Values,omitempty"`
}

type Dependency struct {
	Name         string `yaml:"Name"`
	OriginalPath string `yaml:"OriginalPath"`
	StateHash    string `yaml:"StateHash"`
	ImportHash   string `yaml:"ImportHash"`
}

type Dependencies struct {
	StateHash    string                 `yaml:"StateHash"`
	Dependencies map[string]*Dependency `yaml:"Dependencies"`
}

type Inputs struct {
	StateHash string                 `yaml:"StateHash"`
	Inputs    map[string]*Dependency `yaml:"Inputs"`
}

// Metadata

type MetadataServers struct {
	Url string `yaml:"Url"`
}

type Metadata struct {
	Description string                      `yaml:"Description"`
	Servers     map[string]*MetadataServers `yaml:"Servers"`
}

// Relationship

type Relationship struct {
	Name          string `yaml:"Name"`
	RootNode      string `yaml:"RootNode"`
	OriginalPath  string `yaml:"OriginalPath"`
	StateHash     string `yaml:"StateHash"`
	Uri           string `yaml:"Uri"`
	Version       string `yaml:"Version"`
	IsSameProject bool   `yaml:"IsSameProject"`
}

type Relationships struct {
	StateHash     string                   `yaml:"StateHash"`
	Relationships map[string]*Relationship `yaml:"Relationships"`
}

// Imports

type ImportImport struct {
	Alias *string `yaml:"Alias"`
	Path  string  `yaml:"Path"`
}

type Import struct {
	Name         string        `yaml:"Name"`
	RootNode     string        `yaml:"RootNode"`
	OriginalPath string        `yaml:"OriginalPath"`
	StateHash    string        `yaml:"StateHash"`
	Import       *ImportImport `yaml:"Import"`
	Type         string        `yaml:"Type"`
}

type Imports struct {
	StateHash string             `yaml:"StateHash"`
	Imports   map[string]*Import `yaml:"Imports"`
}

// Auth

type Auth struct {
	Name             string  `yaml:"Name"`
	RootNode         string  `yaml:"RootNode"`
	OriginalPath     string  `yaml:"OriginalPath"`
	StateHash        string  `yaml:"StateHash"`
	Description      *string `yaml:"Description"`
	Scheme           string  `yaml:"Scheme"`
	Format           *string `yaml:"Format"`
	ApplyToAllRoutes bool    `yaml:"ApplyToAllRoutes"`
}

type Auths struct {
	StateHash string           `yaml:"StateHash"`
	Auths     map[string]*Auth `yaml:"Auths"`
}

// Enums

type EnumType string

const (
	EnumType_String EnumType = "String"
	EnumType_Int    EnumType = "Int"
)

type EnumValue struct {
	Name  string
	Value string
}

type Enum struct {
	Name         string       `yaml:"Name"`
	RootNode     string       `yaml:"RootNode"`
	OriginalPath string       `yaml:"OriginalPath"`
	StateHash    string       `yaml:"StateHash"`
	Type         EnumType     `yaml:"Type"`
	Values       []*EnumValue `yaml:"Values"`
}

type Enums struct {
	StateHash string           `yaml:"StateHash"`
	Enums     map[string]*Enum `yaml:"Enums"`
}

// Fields

type Fields struct {
	StateHash string            `yaml:"StateHash"`
	Fields    map[string]*Field `yaml:"Fields"`
}

// Types

type Type struct {
	Name         string   `yaml:"Name"`
	RootNode     string   `yaml:"RootNode"`
	OriginalPath string   `yaml:"OriginalPath"`
	StateHash    string   `yaml:"StateHash"`
	Fields       []string `yaml:"Fields"`
}

type Types struct {
	StateHash string           `yaml:"StateHash"`
	Types     map[string]*Type `yaml:"Types"`
}

// Entities

type EntitiesMetatada struct {
	Schema      string `yaml:"Schema"`
	ColumnsCase string `yaml:"ColumnsCase"`
}

type EntityColumn struct {
	Name         string `yaml:"Name"`
	ColumnName   string `yaml:"ColumnName"`
	OriginalPath string `yaml:"OriginalPath"`
	StateHash    string `yaml:"StateHash"`
	FieldHash    string `yaml:"FieldHash"`
	Type         string `yaml:"Type"`
}

type EntityPrimaryKey struct {
	Hash           string   `yaml:"Hash"`
	ConstraintName string   `yaml:"ConstraintName"`
	ColumnsHashes  []string `yaml:"ColumnsHashes"`
}

type EntityIndex struct {
	ConstraintName string   `yaml:"ConstraintName"`
	ColumnsHashes  []string `yaml:"ColumnsHashes"`
	Unique         bool     `yaml:"Unique"`
}

type EntityForeignKey struct {
	ConstraintName   string   `yaml:"ConstraintName"`
	ColumnsHashes    []string `yaml:"ColumnsHashes"`
	RefTableHash     string   `yaml:"RefTableHash"`
	RefColumnsHashes []string `yaml:"RefColumnsHashes"`
	OnDelete         *string  `yaml:"OnDelete"`
	OnUpdate         *string  `yaml:"OnUpdate"`
}

type Entity struct {
	Name         string                       `yaml:"Name"`
	RootNode     string                       `yaml:"RootNode"`
	OriginalPath string                       `yaml:"OriginalPath"`
	Schema       string                       `yaml:"Schema"`
	TableName    string                       `yaml:"TableName"`
	StateHash    string                       `yaml:"StateHash"`
	Columns      map[string]*EntityColumn     `yaml:"Columns"`
	PrimaryKeys  *EntityPrimaryKey            `yaml:"PrimaryKeys"`
	Indexes      map[string]*EntityIndex      `yaml:"Indexes"`
	ForeignKeys  map[string]*EntityForeignKey `yaml:"ForeignKeys"`
}

type Entities struct {
	StateHash string             `yaml:"StateHash"`
	Metadata  *EntitiesMetatada  `yaml:"Metadata"`
	Entities  map[string]*Entity `yaml:"Entities"`
}

// Repository

type RepositoryMethodInput struct {
	TypeHash string `yaml:"TypeHash"`
}

type RepositoryMethodOutput struct {
	TypeHash string `yaml:"TypeHash"`
}

type RepositoryMethod struct {
	Name         string                  `yaml:"Name"`
	OriginalPath string                  `yaml:"OriginalPath"`
	StateHash    string                  `yaml:"StateHash"`
	Input        *RepositoryMethodInput  `yaml:"Input"`
	Output       *RepositoryMethodOutput `yaml:"Output"`
}

type RepositoryMethods struct {
	StateHash string                       `yaml:"StateHash"`
	Methods   map[string]*RepositoryMethod `yaml:"Methods"`
}

type Repository struct {
	StateHash    string             `yaml:"StateHash"`
	Dependencies *Dependencies      `yaml:"Dependencies"`
	Inputs       *Inputs            `yaml:"Inputs"`
	Methods      *RepositoryMethods `yaml:"Methods"`
}

// Usecase

type UsecaseMethodInput struct {
	TypeHash string `yaml:"TypeHash"`
}

type UsecaseMethodOutput struct {
	TypeHash string `yaml:"TypeHash"`
}

type UsecaseMethod struct {
	Name         string               `yaml:"Name"`
	OriginalPath string               `yaml:"OriginalPath"`
	StateHash    string               `yaml:"StateHash"`
	Input        *UsecaseMethodInput  `yaml:"Input"`
	Output       *UsecaseMethodOutput `yaml:"Output"`
	EventHashes  []string             `yaml:"EventHashes"`
}

type UsecaseMethods struct {
	StateHash string                    `yaml:"StateHash"`
	Methods   map[string]*UsecaseMethod `yaml:"Methods"`
}

type Usecase struct {
	StateHash    string          `yaml:"StateHash"`
	Dependencies *Dependencies   `yaml:"Dependencies"`
	Inputs       *Inputs         `yaml:"Inputs"`
	Methods      *UsecaseMethods `yaml:"Methods"`
}

// Delivery

type DeliveryGrpcRpcExample struct {
	Name         string `yaml:"Name"`
	OriginalPath string `yaml:"OriginalPath"`
	StateHash    string `yaml:"StateHash"`
	StatusCode   int    `yaml:"StatusCode"`
	Message      any    `yaml:"Message"`
	Returns      any    `yaml:"Returns"`
}

type DeliveryGrpcRpc struct {
	MethodHash   string                             `yaml:"MethodHash"`
	OriginalPath string                             `yaml:"OriginalPath"`
	Examples     map[string]*DeliveryGrpcRpcExample `yaml:"Examples"`
}

type DeliveryGrpc struct {
	StateHash string                      `yaml:"StateHash"`
	Rpcs      map[string]*DeliveryGrpcRpc `yaml:"Rpcs"`
}

type Delivery struct {
	StateHash    string        `yaml:"StateHash"`
	Dependencies *Dependencies `yaml:"Dependencies"`
	Grpc         *DeliveryGrpc `yaml:"Grpc"`
}

// Schema

type Schema struct {
	Domain        string         `yaml:"Domain"`
	Version       string         `yaml:"Version"`
	Metadata      *Metadata      `yaml:"Metadata"`
	Relationships *Relationships `yaml:"Relationships"`
	Imports       *Imports       `yaml:"Imports"`
	Auths         *Auths         `yaml:"Auth"`
	Enums         *Enums         `yaml:"Enums"`
	Fields        *Fields        `yaml:"Fields"`
	Types         *Types         `yaml:"Types"`
	Entities      *Entities      `yaml:"Entities"`
	Repository    *Repository    `yaml:"Repository"`
	Usecase       *Usecase       `yaml:"Usecase"`
	Delivery      *Delivery      `yaml:"Delivery"`
}
