package types

type DeliveryGrpc struct {
	GenClient     *bool
	ClientPath    *string
	GenProto      *bool
	ProtofilePath *string
	GenServer     *string
	ServerPath    *string
}

type Delivery struct {
	Dependencies map[string]*Dependency
	Grpc         *DeliveryGrpc
}
