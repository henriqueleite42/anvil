package parser

import (
	"fmt"

	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
)

func (self *Parser) ResolveMethod(kind Kind, name string, inputTypeHash string, outputTypeHash string) error {
	var inputTypeName string
	if inputTypeHash != "" {
		t, ok := self.Schema.Types.Types[inputTypeHash]
		if !ok {
			return fmt.Errorf("type \"%s\" not found for input of method \"%s\"", inputTypeHash, name)
		}

		templT, err := self.ResolveMap(kind, t, "")
		if err != nil {
			return err
		}

		inputTypeName = templT.Name
	}

	var outputTypeName string
	if outputTypeHash != "" {
		t, ok := self.Schema.Types.Types[outputTypeHash]
		if !ok {
			return fmt.Errorf("type \"%s\" not found for output of method \"%s\"", outputTypeHash, name)
		}

		templT, err := self.ResolveMap(kind, t, "")
		if err != nil {
			return err
		}

		outputTypeName = templT.Name
	}

	result := &templates.TemplMethod{
		MethodName:     name,
		InputTypeName:  inputTypeName,
		OutputTypeName: outputTypeName,
	}

	if kind == Kind_Repository {
		self.MethodsRepository = append(self.MethodsRepository, result)
	} else if kind == Kind_Usecase {
		self.MethodsUsecase = append(self.MethodsUsecase, result)
	} else {
		return fmt.Errorf("kind \"%s\" not supported by \"ResolveMethod\"", kind)
	}

	return nil
}
