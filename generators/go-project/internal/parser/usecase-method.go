package parser

import (
	"fmt"

	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *Parser) ResolveUsecaseMethod(usc *schemas.UsecaseMethod, pkgName string) error {
	_, ok := self.MethodsUsecaseToAvoidDuplication[usc.Ref]
	if ok {
		return nil
	}

	var inputTypeName string
	if usc.Input != nil && usc.Input.TypeHash != "" {
		t, ok := self.Schema.Types.Types[usc.Input.TypeHash]
		if !ok {
			return fmt.Errorf("fail to find type for \"%s.Input\"", usc.Name)
		}

		tParsed, err := self.GoTypesParserUsecase.ParseType(t)
		if err != nil {
			return err
		}

		inputTypeName = tParsed.GetFullTypeName(pkgName)
	}

	var outputTypeName string
	if usc.Output != nil && usc.Output.TypeHash != "" {
		t, ok := self.Schema.Types.Types[usc.Output.TypeHash]
		if !ok {
			return fmt.Errorf("fail to find type for \"%s.Output\"", usc.Name)
		}

		tParsed, err := self.GoTypesParserUsecase.ParseType(t)
		if err != nil {
			return err
		}

		outputTypeName = tParsed.GetFullTypeName(pkgName)
	}

	self.MethodsUsecase = append(self.MethodsUsecase, &templates.TemplMethod{
		MethodName:     usc.Name,
		InputTypeName:  inputTypeName,
		OutputTypeName: outputTypeName,
		Order:          usc.Order,
	})

	self.MethodsUsecaseToAvoidDuplication[usc.Ref] = true

	return nil
}
