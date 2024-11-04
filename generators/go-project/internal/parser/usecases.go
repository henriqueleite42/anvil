package parser

import (
	"fmt"

	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/formatter"
	"github.com/henriqueleite42/anvil/language-helpers/golang/imports"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *Parser) resolveUsecaseMethod(usc *schemas.UsecaseMethod) error {
	pkgName := formatter.PascalToSnake(usc.Domain) + "_usecase"

	methodImportsManager := imports.NewImportsManager()

	var inputTypeName string
	if usc.Input != nil && usc.Input.TypeHash != "" {
		t, ok := self.schema.Types.Types[usc.Input.TypeHash]
		if !ok {
			return fmt.Errorf("fail to find type for \"%s.Input\"", usc.Name)
		}

		tParsed, err := self.GoTypesParser.ParseType(t)
		if err != nil {
			return err
		}

		methodImportsManager.MergeImport(tParsed.ModuleImport)

		inputTypeName = tParsed.GetFullTypeName(pkgName)
	}

	var outputTypeName string
	if usc.Output != nil && usc.Output.TypeHash != "" {
		t, ok := self.schema.Types.Types[usc.Output.TypeHash]
		if !ok {
			return fmt.Errorf("fail to find type for \"%s.Output\"", usc.Name)
		}

		tParsed, err := self.GoTypesParser.ParseType(t)
		if err != nil {
			return err
		}

		methodImportsManager.MergeImport(tParsed.ModuleImport)

		outputTypeName = tParsed.GetFullTypeName(pkgName)
	}

	methodImportsManager.AddImport("context", nil)
	methodImportsUnorganized := methodImportsManager.GetImportsUnorganized()
	imports := imports.ResolveImports(methodImportsUnorganized, pkgName)

	self.ImportsUsecase[usc.Domain].MergeImports(methodImportsUnorganized)

	self.usecases[usc.Domain].Methods = append(self.usecases[usc.Domain].Methods, &templates.TemplMethod{
		MethodName:     usc.Name,
		InputTypeName:  inputTypeName,
		OutputTypeName: outputTypeName,
		Order:          usc.Order,
		Imports:        imports,
	})

	return nil
}

func (self *Parser) parseUsecases() error {
	if self.schema.Usecases == nil || self.schema.Usecases.Usecases == nil {
		return nil
	}

	for _, usecase := range self.schema.Usecases.Usecases {
		if usecase.Methods == nil || usecase.Methods.Methods == nil {
			continue
		}

		for _, v := range usecase.Methods.Methods {
			if _, ok := self.usecases[v.Domain]; !ok {
				self.usecases[v.Domain] = &ParserUsecase{
					Methods: make([]*templates.TemplMethod, 0, len(usecase.Methods.Methods)),
				}

				self.ImportsUsecase[v.Domain].AddImport("context", nil)
			}

			err := self.resolveUsecaseMethod(v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
