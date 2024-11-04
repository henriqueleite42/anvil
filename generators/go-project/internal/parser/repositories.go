package parser

import (
	"fmt"

	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/formatter"
	"github.com/henriqueleite42/anvil/language-helpers/golang/imports"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *Parser) resolveRepositoryMethod(rpt *schemas.RepositoryMethod) error {
	pkgName := formatter.PascalToSnake(rpt.Domain) + "_repository"

	methodImportsManager := imports.NewImportsManager()

	var inputTypeName string
	if rpt.Input != nil && rpt.Input.TypeHash != "" {
		t, ok := self.schema.Types.Types[rpt.Input.TypeHash]
		if !ok {
			return fmt.Errorf("fail to find type for \"%s.Input\"", rpt.Name)
		}

		tParsed, err := self.GoTypesParser.ParseType(t)
		if err != nil {
			return err
		}

		methodImportsManager.MergeImport(tParsed.ModuleImport)

		inputTypeName = tParsed.GetFullTypeName(pkgName)
	}

	var outputTypeName string
	if rpt.Output != nil && rpt.Output.TypeHash != "" {
		t, ok := self.schema.Types.Types[rpt.Output.TypeHash]
		if !ok {
			return fmt.Errorf("fail to find type for \"%s.Output\"", rpt.Name)
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

	self.ImportsRepository[rpt.Domain].MergeImports(methodImportsUnorganized)

	self.repositories[rpt.Domain].Methods = append(self.repositories[rpt.Domain].Methods, &templates.TemplMethod{
		MethodName:     rpt.Name,
		InputTypeName:  inputTypeName,
		OutputTypeName: outputTypeName,
		Order:          rpt.Order,
		Imports:        imports,
	})

	return nil
}

func (self *Parser) parseRepositories() error {
	if self.schema.Repositories == nil || self.schema.Repositories.Repositories == nil {
		return nil
	}

	for _, repository := range self.schema.Repositories.Repositories {
		if repository.Methods == nil || repository.Methods.Methods == nil {
			continue
		}

		for _, v := range repository.Methods.Methods {
			if _, ok := self.repositories[v.Domain]; !ok {
				self.repositories[v.Domain] = &ParserRepository{
					Methods: make([]*templates.TemplMethod, 0, len(repository.Methods.Methods)),
				}

				self.ImportsRepository[v.Domain].AddImport("context", nil)
			}

			err := self.resolveRepositoryMethod(v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
