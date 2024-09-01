package schemas

import "slices"

type Dependency struct {
	OriginalPath string `yaml:"OriginalPath"`
	Name         string `yaml:"Name"`
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
	Ref           string  `yaml:"Ref"`
	OriginalPath  string  `yaml:"OriginalPath"`
	Name          string  `yaml:"Name"`
	RootNode      string  `yaml:"RootNode"`
	StateHash     string  `yaml:"StateHash"`
	Uri           string  `yaml:"Uri"`
	Version       string  `yaml:"Version"`
	IsSameProject bool    `yaml:"IsSameProject"`
	Schema        *Schema `yaml:"Schema"`
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
	OriginalPath string        `yaml:"OriginalPath"`
	Name         string        `yaml:"Name"`
	RootNode     string        `yaml:"RootNode"`
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
	Ref              string  `yaml:"Ref"`
	OriginalPath     string  `yaml:"OriginalPath"`
	Name             string  `yaml:"Name"`
	RootNode         string  `yaml:"RootNode"`
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

var (
	EnumType_String EnumType = "String"
	EnumType_Int    EnumType = "Int"
)

type EnumValue struct {
	Name  *string `yaml:"Name,omitempty" json:"Name,omitempty"`
	Value string  `yaml:"Value"`
}

type Enum struct {
	Ref          string       `yaml:"Ref"`
	OriginalPath string       `yaml:"OriginalPath"`
	Name         string       `yaml:"Name"`
	DbType       string       `yaml:"DbType"`
	RootNode     string       `yaml:"RootNode"`
	StateHash    string       `yaml:"StateHash"`
	Type         EnumType     `yaml:"Type"`
	Values       []*EnumValue `yaml:"Values"`
}

type Enums struct {
	StateHash string           `yaml:"StateHash"`
	Enums     map[string]*Enum `yaml:"Enums"`
}

// Types

type TypeType string

var (
	TypeType_String    TypeType = "String"
	TypeType_Int       TypeType = "Int"
	TypeType_Timestamp TypeType = "Timestamp"
	TypeType_Enum      TypeType = "Enum"
	TypeType_Map       TypeType = "Map"
	TypeType_List      TypeType = "List"
)

var TypeTypeArr = []TypeType{
	TypeType_String,
	TypeType_Int,
	TypeType_Timestamp,
	TypeType_Enum,
	TypeType_Map,
	TypeType_List,
}

func ToTypeType(i string) (TypeType, bool) {
	ft := TypeType(i)

	return ft, slices.Contains(TypeTypeArr, ft)
}

type TypeConfidentiality string

var (
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
	Ref             string              `yaml:"Ref"`
	OriginalPath    string              `yaml:"OriginalPath"`
	Name            string              `yaml:"Name"`
	RootNode        string              `yaml:"RootNode"`
	StateHash       string              `yaml:"StateHash"`
	Confidentiality TypeConfidentiality `yaml:"Confidentiality"`
	Optional        bool                `yaml:"Optional"`
	Format          *string             `yaml:"Format,omitempty" json:"Format,omitempty"`
	Validate        []string            `yaml:"Validate,omitempty" json:"Validate,omitempty"`
	Type            TypeType            `yaml:"Type"`
	DbType          *string             `yaml:"DbType,omitempty" json:"DbType,omitempty"`
	// Used for Map and List (List will always only have 1 item inside the slice)
	ChildTypesHashes []string `yaml:"ChildTypesHashes,omitempty" json:"ChildTypesHashes,omitempty"`
	// Used for Enum
	EnumHash *string `yaml:"EnumHash,omitempty" json:"EnumHash,omitempty"`
}

type Types struct {
	StateHash string           `yaml:"StateHash"`
	Types     map[string]*Type `yaml:"Types"`
}

// Events

type Event struct {
	Ref          string   `yaml:"Ref"`
	OriginalPath string   `yaml:"OriginalPath"`
	Name         string   `yaml:"Name"`
	RootNode     string   `yaml:"RootNode"`
	StateHash    string   `yaml:"StateHash"`
	Formats      []string `yaml:"Formats"`
	TypeHash     string   `yaml:"TypeHash"`
}

type Events struct {
	StateHash string            `yaml:"StateHash"`
	Events    map[string]*Event `yaml:"Events"`
}

// Entities

type ColumnsCase string

var (
	ColumnsCase_Snake  ColumnsCase = "snake"
	ColumnsCase_Pascal ColumnsCase = "pascal"
	ColumnsCase_Camel  ColumnsCase = "camel"
)

var ColumnsCaseArr = []ColumnsCase{
	ColumnsCase_Snake,
	ColumnsCase_Pascal,
	ColumnsCase_Camel,
}

func ToColumnsCase(i string) (ColumnsCase, bool) {
	ft := ColumnsCase(i)

	return ft, slices.Contains(ColumnsCaseArr, ft)
}

type EntitiesMetadata struct {
	ColumnsCase *ColumnsCase `yaml:"ColumnsCase,omitempty" json:"ColumnsCase,omitempty"`
}

type EntityColumn struct {
	Ref          string `yaml:"Ref"`
	OriginalPath string `yaml:"OriginalPath"`
	Name         string `yaml:"Name"`
	ColumnName   string `yaml:"ColumnName"`
	StateHash    string `yaml:"StateHash"`
	TypeHash     string `yaml:"TypeHash"`
}

type EntityPrimaryKey struct {
	StateHash      string   `yaml:"StateHash"`
	ConstraintName string   `yaml:"ConstraintName"`
	ColumnsHashes  []string `yaml:"ColumnsHashes"`
}

type EntityIndex struct {
	OriginalPath   string   `yaml:"OriginalPath"`
	StateHash      string   `yaml:"StateHash"`
	ConstraintName string   `yaml:"ConstraintName"`
	ColumnsHashes  []string `yaml:"ColumnsHashes"`
	Unique         bool     `yaml:"Unique"`
}

type EntityForeignKey struct {
	OriginalPath     string   `yaml:"OriginalPath"`
	StateHash        string   `yaml:"StateHash"`
	ConstraintName   string   `yaml:"ConstraintName"`
	ColumnsHashes    []string `yaml:"ColumnsHashes"`
	RefTableHash     string   `yaml:"RefTableHash"`
	RefColumnsHashes []string `yaml:"RefColumnsHashes"`
	OnDelete         *string  `yaml:"OnDelete"`
	OnUpdate         *string  `yaml:"OnUpdate"`
}

type Entity struct {
	Ref          string                       `yaml:"Ref"`
	OriginalPath string                       `yaml:"OriginalPath"`
	Name         string                       `yaml:"Name"`
	RootNode     string                       `yaml:"RootNode"`
	TypeHash     string                       `yaml:"TypeHash"`
	Schema       *string                      `yaml:"Schema,omitempty" json:"Schema,omitempty"`
	TableName    string                       `yaml:"TableName"`
	StateHash    string                       `yaml:"StateHash"`
	Columns      map[string]*EntityColumn     `yaml:"Columns"`
	PrimaryKey   *EntityPrimaryKey            `yaml:"PrimaryKey"`
	Indexes      map[string]*EntityIndex      `yaml:"Indexes,omitempty" json:"Indexes,omitempty"`
	ForeignKeys  map[string]*EntityForeignKey `yaml:"ForeignKeys,omitempty" json:"ForeignKeys,omitempty"`
}

type Entities struct {
	StateHash string             `yaml:"StateHash"`
	Metadata  *EntitiesMetadata  `yaml:"Metadata"`
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
	Ref          string                  `yaml:"Ref"`
	OriginalPath string                  `yaml:"OriginalPath"`
	StateHash    string                  `yaml:"StateHash"`
	Name         string                  `yaml:"Name"`
	Description  *string                 `yaml:"Description,omitempty" json:"Description,omitempty"`
	Input        *RepositoryMethodInput  `yaml:"Input,omitempty" json:"Input,omitempty"`
	Output       *RepositoryMethodOutput `yaml:"Output,omitempty" json:"Output,omitempty"`
}

type RepositoryMethods struct {
	StateHash string                       `yaml:"StateHash"`
	Methods   map[string]*RepositoryMethod `yaml:"Methods"`
}

type Repository struct {
	StateHash    string             `yaml:"StateHash"`
	Dependencies *Dependencies      `yaml:"Dependencies,omitempty" json:"Dependencies,omitempty"`
	Inputs       *Inputs            `yaml:"Inputs,omitempty" json:"Inputs,omitempty"`
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
	Ref          string               `yaml:"Ref"`
	OriginalPath string               `yaml:"OriginalPath"`
	Name         string               `yaml:"Name"`
	Description  *string              `yaml:"Description,omitempty" json:"Description,omitempty"`
	StateHash    string               `yaml:"StateHash"`
	Input        *UsecaseMethodInput  `yaml:"Input,omitempty" json:"Input,omitempty"`
	Output       *UsecaseMethodOutput `yaml:"Output,omitempty" json:"Output,omitempty"`
	EventHashes  []string             `yaml:"EventHashes,omitempty" json:"EventHashes,omitempty"`
}

type UsecaseMethods struct {
	StateHash string                    `yaml:"StateHash"`
	Methods   map[string]*UsecaseMethod `yaml:"Methods"`
}

type Usecase struct {
	StateHash    string          `yaml:"StateHash"`
	Dependencies *Dependencies   `yaml:"Dependencies,omitempty" json:"Dependencies,omitempty"`
	Inputs       *Inputs         `yaml:"Inputs,omitempty" json:"Inputs,omitempty"`
	Methods      *UsecaseMethods `yaml:"Methods,omitempty"`
}

// Delivery

type DeliveryGrpcRpcExample struct {
	OriginalPath string `yaml:"OriginalPath"`
	Name         string `yaml:"Name"`
	StateHash    string `yaml:"StateHash"`
	StatusCode   int    `yaml:"StatusCode"`
	Message      any    `yaml:"Message"`
	Returns      any    `yaml:"Returns"`
}

type DeliveryGrpcRpc struct {
	Ref               string                             `yaml:"Ref"`
	OriginalPath      string                             `yaml:"OriginalPath"`
	UsecaseMethodHash string                             `yaml:"UsecaseMethodHash"`
	Examples          map[string]*DeliveryGrpcRpcExample `yaml:"Examples"`
}

type DeliveryGrpc struct {
	StateHash string                      `yaml:"StateHash"`
	Rpcs      map[string]*DeliveryGrpcRpc `yaml:"Rpcs"`
}

type DeliveryHttpRouteExample struct {
	OriginalPath string `yaml:"OriginalPath"`
	Name         string `yaml:"Name"`
	StateHash    string `yaml:"StateHash"`
	// TODO
}

type DeliveryHttpRoute struct {
	Ref                string                               `yaml:"Ref"`
	OriginalPath       string                               `yaml:"OriginalPath"`
	UsecaseMethodHash  string                               `yaml:"UsecaseMethodHash"`
	StatusCode         int                                  `yaml:"StatusCode"`
	HttpMethod         string                               `yaml:"HttpMethod"`
	Path               string                               `yaml:"Path"`
	ReqHeadersTypeHash *string                              `yaml:"ReqHeadersTypesHashes,omitempty" json:"ReqHeadersTypesHashes,omitempty"`
	ResHeadersTypeHash *string                              `yaml:"ResHeadersTypesHashes,omitempty" json:"ResHeadersTypesHashes,omitempty"`
	Auth               *string                              `yaml:"Auth,omitempty" json:"Auth,omitempty"`
	Examples           map[string]*DeliveryHttpRouteExample `yaml:"Examples"`
}

type DeliveryHttp struct {
	StateHash string                        `yaml:"StateHash"`
	Routes    map[string]*DeliveryHttpRoute `yaml:"Routes"`
}

type DeliveryQueueQueueExample struct {
	OriginalPath string `yaml:"OriginalPath"`
	Name         string `yaml:"Name"`
	StateHash    string `yaml:"StateHash"`
	// TODO
}

type DeliveryQueueQueue struct {
	Ref               string                                `yaml:"Ref"`
	OriginalPath      string                                `yaml:"OriginalPath"`
	UsecaseMethodHash string                                `yaml:"UsecaseMethodHash"`
	QueueId           string                                `yaml:"QueueId"`
	Bulk              bool                                  `yaml:"Bulk"`
	Examples          map[string]*DeliveryQueueQueueExample `yaml:"Examples"`
}

type DeliveryQueue struct {
	StateHash string                         `yaml:"StateHash"`
	Queues    map[string]*DeliveryQueueQueue `yaml:"Queues"`
}

type Delivery struct {
	StateHash    string         `yaml:"StateHash"`
	Dependencies *Dependencies  `yaml:"Dependencies"`
	Grpc         *DeliveryGrpc  `yaml:"Grpc"`
	Http         *DeliveryHttp  `yaml:"Http"`
	Queue        *DeliveryQueue `yaml:"Queue"`
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
