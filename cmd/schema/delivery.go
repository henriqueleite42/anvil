package schema

type DeliveryGrpc struct {
	GenClient     *bool   `yaml:"GenClient,omitempty"`
	ClientPath    *string `yaml:"ClientPath,omitempty"`
	GenProto      *bool   `yaml:"GenProto,omitempty"`
	ProtofilePath *string `yaml:"ProtofilePath,omitempty"`
	GenServer     *string `yaml:"GenServer,omitempty"`
	ServerPath    *string `yaml:"ServerPath,omitempty"`
}

type Delivery struct {
	Dependencies map[string]*Dependency `yaml:"Dependencies,omitempty"`
	Grpc         *DeliveryGrpc          `yaml:"Grpc,omitempty"`
}
