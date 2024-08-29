package schema

import "slices"

// Common

type TypeType string

const (
	TypeType_String        TypeType = "String"
	TypeType_Int           TypeType = "Int"
	TypeType_Timestamp     TypeType = "Timestamp"
	TypeType_Enum          TypeType = "Enum"
	TypeType_Map           TypeType = "Map"
	TypeType_MapStringMap  TypeType = "Map[String]Map"
	TypeType_ListString    TypeType = "List[String]"
	TypeType_ListInt       TypeType = "List[Int]"
	TypeType_ListTimestamp TypeType = "List[Timestamp]"
	TypeType_ListEnum      TypeType = "List[Enum]"
	TypeType_ListMap       TypeType = "List[Map]"
)

var TypeTypeArr = []TypeType{
	TypeType_String,
	TypeType_Int,
	TypeType_Timestamp,
	TypeType_Enum,
	TypeType_Map,
	TypeType_MapStringMap,
	TypeType_ListString,
	TypeType_ListInt,
	TypeType_ListTimestamp,
	TypeType_ListEnum,
	TypeType_ListMap,
}

func ToTypeType(i string) (TypeType, bool) {
	ft := TypeType(i)

	return ft, slices.Contains(TypeTypeArr, ft)
}

type TypeConfidentiality string

const (
	TypeConfidentiality_Low    TypeConfidentiality = "LOW"
	TypeConfidentiality_Medium TypeConfidentiality = "MEDIUM"
	TypeConfidentiality_High   TypeConfidentiality = "HIGH"
)

var TypeConfidentialityArr = []TypeConfidentiality{
	TypeConfidentiality_Low,
	TypeConfidentiality_Medium,
	TypeConfidentiality_High,
}

func ToTypeConfidentiality(i string) (TypeConfidentiality, bool) {
	ft := TypeConfidentiality(i)

	return ft, slices.Contains(TypeConfidentialityArr, ft)
}

type Type struct {
	Name            string              `yaml:"Name"`
	RootNode        string              `yaml:"RootNode"`
	OriginalPath    string              `yaml:"OriginalPath"`
	StateHash       string              `yaml:"StateHash"`
	Confidentiality TypeConfidentiality `yaml:"Confidentiality"`
	Optional        bool                `yaml:"Optional"`
	Format          *string             `yaml:"Format,omitempty" json:"Format,omitempty"`
	Validate        []string            `yaml:"Validate,omitempty" json:"Validate,omitempty"`
	Type            TypeType            `yaml:"Type"`
	DbType          *string             `yaml:"DbType,omitempty" json:"DbType,omitempty"`
	// Used for Map and List[Map]
	ChildTypesHashes []string `yaml:"ChildTypesHashes,omitempty" json:"ChildTypesHashes,omitempty"`
	// Used for Enum and List[Enum]
	EnumHash *string `yaml:"EnumHash,omitempty" json:"EnumHash,omitempty"`
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
	Name  *string `yaml:"Name,omitempty" json:"Name,omitempty"`
	Value string  `yaml:"Value"`
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

// Types

type Types struct {
	StateHash string           `yaml:"StateHash"`
	Types     map[string]*Type `yaml:"Types"`
}

// Events

type Event struct {
	Name         string   `yaml:"Name"`
	RootNode     string   `yaml:"RootNode"`
	OriginalPath string   `yaml:"OriginalPath"`
	StateHash    string   `yaml:"StateHash"`
	Formats      []string `yaml:"Formats"`
	TypeHash     string   `yaml:"TypeHash"`
}

type Events struct {
	StateHash string            `yaml:"StateHash"`
	Events    map[string]*Event `yaml:"Events"`
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
	TypeHash     string `yaml:"TypeHash"`
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
	Domain        string         `yaml:"Domain,omitempty" json:"Domain,omitempty"`
	Version       string         `yaml:"Version,omitempty" json:"Version,omitempty"`
	Metadata      *Metadata      `yaml:"Metadata,omitempty" json:"Metadata,omitempty"`
	Relationships *Relationships `yaml:"Relationships,omitempty" json:"Relationships,omitempty"`
	Imports       *Imports       `yaml:"Imports,omitempty" json:"Imports,omitempty"`
	Auths         *Auths         `yaml:"Auth,omitempty" json:"Auth,omitempty"`
	Enums         *Enums         `yaml:"Enums,omitempty" json:"Enums,omitempty"`
	Types         *Types         `yaml:"Types,omitempty" json:"Types,omitempty"`
	Events        *Events        `yaml:"Events,omitempty" json:"Events,omitempty"`
	Entities      *Entities      `yaml:"Entities,omitempty" json:"Entities,omitempty"`
	Repository    *Repository    `yaml:"Repository,omitempty" json:"Repository,omitempty"`
	Usecase       *Usecase       `yaml:"Usecase,omitempty" json:"Usecase,omitempty"`
	Delivery      *Delivery      `yaml:"Delivery,omitempty" json:"Delivery,omitempty"`
}
