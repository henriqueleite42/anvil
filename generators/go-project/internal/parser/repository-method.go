package parser

import (
	"fmt"

	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

func (self *Parser) ResolveRepositoryMethod(usc *schemas.RepositoryMethod) error {
	var inputTypeName string
	if usc.Input != nil && usc.Input.TypeHash != "" {
		t, ok := self.Schema.Types.Types[usc.Input.TypeHash]
		if !ok {
			return fmt.Errorf("fail to find type for \"%s.Input\"", usc.Name)
		}

		tParsed, err := self.GoTypesParserUsecase.ParseType(t, &types_parser.ParseTypeOpt{
			PrefixForEnums: "models",
		})
		if err != nil {
			return err
		}

		inputTypeName = tParsed.GolangType
	}

	var outputTypeName string
	if usc.Output != nil && usc.Output.TypeHash != "" {
		t, ok := self.Schema.Types.Types[usc.Output.TypeHash]
		if !ok {
			return fmt.Errorf("fail to find type for \"%s.Output\"", usc.Name)
		}

		tParsed, err := self.GoTypesParserUsecase.ParseType(t, &types_parser.ParseTypeOpt{
			PrefixForEnums: "models",
		})
		if err != nil {
			return err
		}

		outputTypeName = tParsed.GolangType
	}

	self.MethodsRepository = append(self.MethodsRepository, &templates.TemplMethod{
		MethodName:     usc.Name,
		InputTypeName:  inputTypeName,
		OutputTypeName: outputTypeName,
	})

	return nil
}
