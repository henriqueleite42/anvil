package parser

import (
	"strings"

	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

type ModelsParsed struct {
	Imports  [][]string
	Enums    []*types_parser.Enum
	Entities []*types_parser.Type
	Events   []*types_parser.Type
}

func (self *Parser) parseEnums(curDomain string) error {
	if self.Schema.Enums == nil || self.Schema.Enums.Enums == nil {
		return nil
	}

	for _, v := range self.Schema.Enums.Enums {
		if !strings.HasPrefix(v.Ref, curDomain) {
			continue
		}

		_, err := self.GoTypesParserModels.ParseEnum(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *Parser) parseTypes(curDomain string) error {
	if self.Schema.Types == nil || self.Schema.Types.Types == nil {
		return nil
	}

	for _, v := range self.Schema.Types.Types {
		if !strings.HasPrefix(v.Ref, curDomain) {
			continue
		}

		_, err := self.GoTypesParserModels.ParseType(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *Parser) parseEntities(curDomain string) error {
	if self.Schema.Entities == nil || self.Schema.Entities.Entities == nil {
		return nil
	}

	for _, v := range self.Schema.Entities.Entities {
		if !strings.HasPrefix(v.Ref, curDomain) {
			continue
		}

		entity := self.Schema.Types.Types[v.TypeHash]

		_, err := self.GoTypesParserModels.ParseType(entity)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *Parser) ParseModels(curDomain string) (*ModelsParsed, error) {
	if self.GoTypesParserModels == nil {
		return &ModelsParsed{}, nil
	}

	err := self.parseEnums(curDomain)
	if err != nil {
		return nil, err
	}

	err = self.parseTypes(curDomain)
	if err != nil {
		return nil, err
	}

	err = self.parseEntities(curDomain)
	if err != nil {
		return nil, err
	}

	importsModels := self.GoTypesParserModels.GetImports()
	enums := self.GoTypesParserModels.GetEnums()
	entities := self.GoTypesParserModels.GetEntities()
	entities = append(entities,
		self.GoTypesParserModels.GetTypes()...,
	)
	events := self.GoTypesParserModels.GetEvents()

	return &ModelsParsed{
		Imports:  importsModels,
		Enums:    enums,
		Entities: entities,
		Events:   events,
	}, nil
}
