package schema

type MethodDeliveryGrpc_Example struct {
	StatusCode int `yaml:"StatusCode"`
	Message    any `yaml:"Message,omitempty"`
	Returns    any `yaml:"Returns,omitempty"`
}

type MethodDeliveryGrpc struct {
	Examples map[string]*MethodDeliveryGrpc_Example `yaml:"Examples,omitempty"`
}

type MethodDeliveryQueue struct {
	Id        string `yaml:"Id"`
	RelatedTo string `yaml:"RelatedTo"`
}

type MethodDelivery struct {
	Grpc  *MethodDeliveryGrpc  `yaml:"Grpc,omitempty"`
	Queue *MethodDeliveryQueue `yaml:"Queue,omitempty"`
}

type Method struct {
	Input    map[string]*Field `yaml:"Input,omitempty"`
	Output   map[string]*Field `yaml:"Output,omitempty"`
	Delivery *MethodDelivery   `yaml:"Delivery,omitempty"`
}

type Usecase struct {
	Dependencies map[string]*Dependency `yaml:"Dependencies,omitempty"`
	Inputs       map[string]*Dependency `yaml:"Inputs,omitempty"`
	Methods      map[string]*Method     `yaml:"Methods,omitempty"`
}
