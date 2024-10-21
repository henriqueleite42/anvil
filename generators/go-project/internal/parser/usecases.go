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

type ParsedUsecase struct {
	Imports [][]string
	Types   []*types_parser.Type
	Methods []*templates.TemplMethod
}

func (self *Parser) resolveUsecaseMethod(usc *schemas.UsecaseMethod, pkgName string) error {
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

func (self *Parser) ParseUsecases(curDomain string) (*ParsedUsecase, error) {
	if self.Schema.Usecases == nil || self.Schema.Usecases.Usecases == nil {
		return &ParsedUsecase{}, nil
	}

	Usecase, ok := self.Schema.Usecases.Usecases[curDomain]
	if !ok {
		return &ParsedUsecase{}, nil
	}

	self.MethodsUsecaseToAvoidDuplication = map[string]bool{}
	self.MethodsUsecase = []*templates.TemplMethod{}

	for _, v := range Usecase.Methods.Methods {
		if !strings.HasPrefix(v.Ref, curDomain) {
			continue
		}

		domainSnake := formatter.PascalToSnake(curDomain)

		err := self.resolveUsecaseMethod(v, domainSnake+"_usecase")
		if err != nil {
			return nil, err
		}
	}

	self.GoTypesParserUsecase.AddImport("context")
	imports := self.GoTypesParserUsecase.GetImports()
	self.GoTypesParserUsecase.ResetImports()

	types := self.GoTypesParserUsecase.GetUsecase()
	sort.Slice(types, func(i, j int) bool {
		return types[i].GolangType < types[j].GolangType
	})

	sort.Slice(self.MethodsUsecase, func(i, j int) bool {
		return self.MethodsUsecase[i].Order < self.MethodsUsecase[j].Order
	})

	return &ParsedUsecase{
		Imports: imports,
		Types:   types,
		Methods: self.MethodsUsecase,
	}, nil
}
