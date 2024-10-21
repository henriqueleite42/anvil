package adapters

type Validator interface {
	Validate(i any) error
}
