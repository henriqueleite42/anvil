package templates

const ValidatorTempl = `package adapters

type Validator interface {
	Validate(i interface{}) error
}
`

const ValidatorImplementationTempl = `package go_validator

import (
	"github.com/go-playground/validator/v10"
)

type playgroundValidator struct {
	validator *validator.Validate
}

func (self *playgroundValidator) Validate(i interface{}) error {
	return self.validator.Struct(i)
}

func NewGoValidator() (adapters.Validator, error) {
	vald := validator.New(validator.WithRequiredStructEnabled())

	// Add custom validations here

	return &playgroundValidator{
		validator: vald,
	}, nil
}
`
