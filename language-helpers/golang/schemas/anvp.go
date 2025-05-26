package schemas

// Schema

type Schema struct {
	Domain      string  `yaml:"Domain"`
	Description *string `yaml:"Description,omitempty" json:"Description,omitempty"`
	Version     *string `yaml:"Version,omitempty" json:"Version,omitempty"`
	Uri         string  `yaml:"Uri"`
}

// Auth

type Auth struct {
	Ref          string `yaml:"Ref"`
	OriginalPath string `yaml:"OriginalPath"`
	RootNode     string `yaml:"RootNode"`
	Domain       string `yaml:"Domain"`
	StateHash    string `yaml:"StateHash"`

	Name             string  `yaml:"Name"`
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

type EnumValue struct {
	Name       string `yaml:"Name"`
	Value      string `yaml:"Value"`
	Index      uint   `yaml:"Index"`
	Deprecated bool   `yaml:"Deprecated"`
}

type Enum struct {
	Ref          string `yaml:"Ref"`
	OriginalPath string `yaml:"OriginalPath"`
	RootNode     string `yaml:"RootNode"`
	Domain       string `yaml:"Domain"`
	StateHash    string `yaml:"StateHash"`

	Name   string       `yaml:"Name"`
	DbName string       `yaml:"DbName"`
	DbType string       `yaml:"DbType"`
	Type   EnumType     `yaml:"Type"`
	Values []*EnumValue `yaml:"Values"`
}

type Enums struct {
	StateHash string           `yaml:"StateHash"`
	Enums     map[string]*Enum `yaml:"Enums"`
}

// Types

type TypeChild struct {
	PropName *string `yaml:"PropName,omitempty" json:"PropName,omitempty"` // Only present in Map parent types
	TypeHash string  `yaml:"TypeHash"`
}

type Type struct {
	Ref          string `yaml:"Ref"`
	OriginalPath string `yaml:"OriginalPath"`
	RootNode     string `yaml:"RootNode"`
	Domain       string `yaml:"Domain"`
	StateHash    string `yaml:"StateHash"`

	Name            string              `yaml:"Name"` // This is the TYPE name, not the property name or variable name
	Confidentiality TypeConfidentiality `yaml:"Confidentiality"`
	Optional        bool                `yaml:"Optional"`
	Format          *string             `yaml:"Format,omitempty" json:"Format,omitempty"`
	Validate        []string            `yaml:"Validate,omitempty" json:"Validate,omitempty"`
	Transform       []string            `yaml:"Transform,omitempty" json:"Transform,omitempty"`
	AutoIncrement   bool                `yaml:"AutoIncrement"`
	Default         *string             `yaml:"Default,omitempty" json:"Default,omitempty"`
	Type            TypeType            `yaml:"Type"`
	DbName          *string             `yaml:"DbName,omitempty" json:"DbName,omitempty"`
	DbType          *string             `yaml:"DbType,omitempty" json:"DbType,omitempty"`
	// Used for Map and List (List will always only have 1 item inside the slice)
	ChildTypes []*TypeChild `yaml:"ChildTypes,omitempty" json:"ChildTypes,omitempty"`
	// Used for Enum
	EnumHash *string `yaml:"EnumHash,omitempty" json:"EnumHash,omitempty"`
}

type Types struct {
	StateHash string           `yaml:"StateHash"`
	Types     map[string]*Type `yaml:"Types"`
}

// Events

type Event struct {
	Ref          string `yaml:"Ref"`
	OriginalPath string `yaml:"OriginalPath"`
	RootNode     string `yaml:"RootNode"`
	Domain       string `yaml:"Domain"`
	StateHash    string `yaml:"StateHash"`

	Name     string   `yaml:"Name"`
	Formats  []string `yaml:"Formats"`
	TypeHash string   `yaml:"TypeHash"`
}

type Events struct {
	StateHash string            `yaml:"StateHash"`
	Events    map[string]*Event `yaml:"Events"`
}

// Entities

type EntitiesMetadata struct {
	NamingCase *NamingCase `yaml:"NamingCase,omitempty" json:"NamingCase,omitempty"`
}

type EntityColumn struct {
	Ref          string `yaml:"Ref"`
	OriginalPath string `yaml:"OriginalPath"`
	Order        uint   `yaml:"Order"`

	Name      string `yaml:"Name"`
	StateHash string `yaml:"StateHash"`
	TypeHash  string `yaml:"TypeHash"`
}

type EntityPrimaryKey struct {
	StateHash      string   `yaml:"StateHash"`
	ConstraintName string   `yaml:"ConstraintName"`
	ColumnsHashes  []string `yaml:"ColumnsHashes"`
}

type EntityIndex struct {
	OriginalPath string `yaml:"OriginalPath"`
	StateHash    string `yaml:"StateHash"`
	Order        uint   `yaml:"Order"`

	ConstraintName string   `yaml:"ConstraintName"`
	ColumnsHashes  []string `yaml:"ColumnsHashes"`
	Unique         bool     `yaml:"Unique"`
}

type EntityForeignKey struct {
	OriginalPath string `yaml:"OriginalPath"`
	StateHash    string `yaml:"StateHash"`
	Order        uint   `yaml:"Order"`

	ConstraintName   string   `yaml:"ConstraintName"`
	ColumnsHashes    []string `yaml:"ColumnsHashes"`
	RefTableHash     string   `yaml:"RefTableHash"`
	RefColumnsHashes []string `yaml:"RefColumnsHashes"`
	OnDelete         *string  `yaml:"OnDelete,omitempty" json:"OnDelete,omitempty"`
	OnUpdate         *string  `yaml:"OnUpdate,omitempty" json:"OnUpdate,omitempty"`
}

type Entity struct {
	Ref          string `yaml:"Ref"`
	OriginalPath string `yaml:"OriginalPath"`
	Domain       string `yaml:"Domain"`
	RootNode     string `yaml:"RootNode"`
	StateHash    string `yaml:"StateHash"`
	Order        uint   `yaml:"Order"`

	Name        string                       `yaml:"Name"`
	TypeHash    string                       `yaml:"TypeHash"`
	DbSchema    *string                      `yaml:"DbSchema,omitempty" json:"DbSchema,omitempty"`
	DbName      string                       `yaml:"DbName"`
	Columns     map[string]*EntityColumn     `yaml:"Columns"`
	PrimaryKey  *EntityPrimaryKey            `yaml:"PrimaryKey"`
	Indexes     map[string]*EntityIndex      `yaml:"Indexes,omitempty" json:"Indexes,omitempty"`
	ForeignKeys map[string]*EntityForeignKey `yaml:"ForeignKeys,omitempty" json:"ForeignKeys,omitempty"`
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
	Ref          string `yaml:"Ref"`
	OriginalPath string `yaml:"OriginalPath"`
	Domain       string `yaml:"Domain"`
	StateHash    string `yaml:"StateHash"`
	Order        uint   `yaml:"Order"`

	Name        string                  `yaml:"Name"`
	Description *string                 `yaml:"Description,omitempty" json:"Description,omitempty"`
	Input       *RepositoryMethodInput  `yaml:"Input,omitempty" json:"Input,omitempty"`
	Output      *RepositoryMethodOutput `yaml:"Output,omitempty" json:"Output,omitempty"`
}

type RepositoryMethods struct {
	StateHash string                       `yaml:"StateHash"`
	Methods   map[string]*RepositoryMethod `yaml:"Methods"`
}

type Repository struct {
	StateHash string             `yaml:"StateHash"`
	Methods   *RepositoryMethods `yaml:"Methods"`
}

type Repositories struct {
	StateHash    string                 `yaml:"StateHash"`
	Repositories map[string]*Repository `yaml:"Repositories"` // [Domain]: Repository
}

// Usecase

type UsecaseMethodInput struct {
	TypeHash string `yaml:"TypeHash"`
}

type UsecaseMethodOutput struct {
	TypeHash string `yaml:"TypeHash"`
}

type UsecaseMethod struct {
	Ref          string `yaml:"Ref"`
	OriginalPath string `yaml:"OriginalPath"`
	Domain       string `yaml:"Domain"`
	StateHash    string `yaml:"StateHash"`
	Order        uint   `yaml:"Order"`

	Name        string               `yaml:"Name"`
	Description *string              `yaml:"Description,omitempty" json:"Description,omitempty"`
	Input       *UsecaseMethodInput  `yaml:"Input,omitempty" json:"Input,omitempty"`
	Output      *UsecaseMethodOutput `yaml:"Output,omitempty" json:"Output,omitempty"`
	EventHashes []string             `yaml:"EventHashes,omitempty" json:"EventHashes,omitempty"`
}

type UsecaseMethods struct {
	StateHash string                    `yaml:"StateHash"`
	Methods   map[string]*UsecaseMethod `yaml:"Methods"`
}

type Usecase struct {
	StateHash string          `yaml:"StateHash"`
	Methods   *UsecaseMethods `yaml:"Methods"`
}

type Usecases struct {
	StateHash string              `yaml:"StateHash"`
	Usecases  map[string]*Usecase `yaml:"Usecases"` // [Domain]: Usecase
}

// Delivery

type DeliveryServers struct {
	Url string `yaml:"Url"`
}

type DeliveryGrpcRpcExample struct {
	OriginalPath string `yaml:"OriginalPath"`
	Domain       string `yaml:"Domain"`
	StateHash    string `yaml:"StateHash"`

	Name       string `yaml:"Name"`
	StatusCode uint   `yaml:"StatusCode"`
	Message    any    `yaml:"Message"`
	Returns    any    `yaml:"Returns"`
}

type DeliveryGrpcRpc struct {
	Ref          string `yaml:"Ref"`
	OriginalPath string `yaml:"OriginalPath"`
	Domain       string `yaml:"Domain"`
	StateHash    string `yaml:"StateHash"`
	Order        uint   `yaml:"Order"`

	Name              string                             `yaml:"Name"`
	UsecaseMethodHash string                             `yaml:"UsecaseMethodHash"`
	Examples          map[string]*DeliveryGrpcRpcExample `yaml:"Examples,omitempty" json:"Examples,omitempty"`
}

type DeliveryGrpc struct {
	StateHash string                      `yaml:"StateHash"`
	Rpcs      map[string]*DeliveryGrpcRpc `yaml:"Rpcs"`
}

type DeliveryHttpRouteExample struct {
	OriginalPath string `yaml:"OriginalPath"`
	Domain       string `yaml:"Domain"`
	StateHash    string `yaml:"StateHash"`

	Name string `yaml:"Name"`
	// TODO
}

type DeliveryHttpRoute struct {
	Ref          string `yaml:"Ref"`
	OriginalPath string `yaml:"OriginalPath"`
	Domain       string `yaml:"Domain"`
	StateHash    string `yaml:"StateHash"`
	Order        uint   `yaml:"Order"`

	UsecaseMethodHash  string                               `yaml:"UsecaseMethodHash"`
	StatusCode         uint                                 `yaml:"StatusCode"`
	HttpMethod         string                               `yaml:"HttpMethod"`
	Path               string                               `yaml:"Path"`
	ReqHeadersTypeHash *string                              `yaml:"ReqHeadersTypesHashes,omitempty" json:"ReqHeadersTypesHashes,omitempty"`
	ResHeadersTypeHash *string                              `yaml:"ResHeadersTypesHashes,omitempty" json:"ResHeadersTypesHashes,omitempty"`
	Auth               *string                              `yaml:"Auth,omitempty" json:"Auth,omitempty"`
	ZipRes             *bool                                `yaml:"ZipRes,omitempty" json:"ZipRes,omitempty"`
	Examples           map[string]*DeliveryHttpRouteExample `yaml:"Examples,omitempty" json:"Examples,omitempty"`
}

type DeliveryHttp struct {
	StateHash string                        `yaml:"StateHash"`
	Routes    map[string]*DeliveryHttpRoute `yaml:"Routes"`
}

type DeliveryQueueQueueExample struct {
	OriginalPath string `yaml:"OriginalPath"`
	Domain       string `yaml:"Domain"`
	StateHash    string `yaml:"StateHash"`
	Order        uint   `yaml:"Order"`

	Name string `yaml:"Name"`
	// TODO
}

type DeliveryQueueQueue struct {
	Ref          string `yaml:"Ref"`
	OriginalPath string `yaml:"OriginalPath"`
	Domain       string `yaml:"Domain"`
	StateHash    string `yaml:"StateHash"`

	UsecaseMethodHash string                                `yaml:"UsecaseMethodHash"`
	QueueId           string                                `yaml:"QueueId"`
	Bulk              bool                                  `yaml:"Bulk"`
	Examples          map[string]*DeliveryQueueQueueExample `yaml:"Examples,omitempty" json:"Examples,omitempty"`
}

type DeliveryQueue struct {
	StateHash string                         `yaml:"StateHash"`
	Queues    map[string]*DeliveryQueueQueue `yaml:"Queues"`
}

type Delivery struct {
	StateHash string                      `yaml:"StateHash"`
	Servers   map[string]*DeliveryServers `yaml:"Servers"`
	Grpc      *DeliveryGrpc               `yaml:"Grpc,omitempty" json:"Grpc,omitempty"`
	Http      *DeliveryHttp               `yaml:"Http,omitempty" json:"Http,omitempty"`
	Queue     *DeliveryQueue              `yaml:"Queue,omitempty" json:"Queue,omitempty"`
}

type Deliveries struct {
	StateHash  string               `yaml:"StateHash"`
	Deliveries map[string]*Delivery `yaml:"Deliveries"` // [Domain]: Delivery
}

// Schema

type AnvpSchema struct {
	// Metadata about domains
	Schemas map[string]*Schema `yaml:"Schemas,omitempty" json:"Schemas,omitempty"`

	// Common to all domains
	Auths    *Auths    `yaml:"Auth,omitempty" json:"Auth,omitempty"`
	Enums    *Enums    `yaml:"Enums,omitempty" json:"Enums,omitempty"`
	Types    *Types    `yaml:"Types,omitempty" json:"Types,omitempty"`
	Events   *Events   `yaml:"Events,omitempty" json:"Events,omitempty"`
	Entities *Entities `yaml:"Entities,omitempty" json:"Entities,omitempty"`

	// Grouped by Domain
	Repositories *Repositories `yaml:"Repository,omitempty" json:"Repository,omitempty"`
	Usecases     *Usecases     `yaml:"Usecase,omitempty" json:"Usecase,omitempty"`
	Deliveries   *Deliveries   `yaml:"Delivery,omitempty" json:"Delivery,omitempty"`
}
