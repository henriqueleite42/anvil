package types

type MethodDeliveryGrpc struct {
	Client *bool
}

type MethodDelivery struct {
	Grpc *MethodDeliveryGrpc
}

type Method struct {
	Input    map[string]*Field
	Output   map[string]*Field
	Delivery *MethodDelivery
}

type Usecase struct {
	Dependencies map[string]*Dependency
	Methods      map[string]*Method
}
