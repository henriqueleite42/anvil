package parser

import (
	"fmt"
	"sort"
	"strings"

	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/formatter"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

type RepositoryParsed struct {
	Imports [][]string
	Types   []*types_parser.Type
	Methods []*templates.TemplMethod
}

func (self *Parser) resolveRepositoryMethod(usc *schemas.RepositoryMethod, pkgName string) error {
	var inputTypeName string
	if usc.Input != nil && usc.Input.TypeHash != "" {
		t, ok := self.Schema.Types.Types[usc.Input.TypeHash]
		if !ok {
			return fmt.Errorf("fail to find type for \"%s.Input\"", usc.Name)
		}

		tParsed, err := self.GoTypesParserRepository.ParseType(t)
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

		tParsed, err := self.GoTypesParserRepository.ParseType(t)
		if err != nil {
			return err
		}

		outputTypeName = tParsed.GetFullTypeName(pkgName)
	}

	self.MethodsRepository = append(self.MethodsRepository, &templates.TemplMethod{
		MethodName:     usc.Name,
		InputTypeName:  inputTypeName,
		OutputTypeName: outputTypeName,
		Order:          usc.Order,
	})

	return nil
}

func (self *Parser) ParseRepositories(curDomain string) (*RepositoryParsed, error) {
	if self.Schema.Repositories == nil || self.Schema.Repositories.Repositories == nil {
		return &RepositoryParsed{}, nil
	}

	repositories, ok := self.Schema.Repositories.Repositories[curDomain]
	if !ok {
		return &RepositoryParsed{}, nil
	}

	self.MethodsRepository = []*templates.TemplMethod{}

	for _, v := range repositories.Methods.Methods {
		if !strings.HasPrefix(v.Ref, curDomain) {
			continue
		}

		domainSnake := formatter.PascalToSnake(curDomain)

		err := self.resolveRepositoryMethod(v, domainSnake+"_repository")
		if err != nil {
			return nil, err
		}
	}

	imports := self.GoTypesParserRepository.GetImports()
	types := self.GoTypesParserRepository.GetTypes()
	sort.Slice(types, func(i, j int) bool {
		return types[i].GolangType < types[j].GolangType
	})
	sort.Slice(self.MethodsRepository, func(i, j int) bool {
		return self.MethodsRepository[i].Order < self.MethodsRepository[j].Order
	})

	return &RepositoryParsed{
		Imports: imports,
		Types:   types,
		Methods: self.MethodsRepository,
	}, nil
}
