package parser

import (
	"github.com/anvil/anvil/internal/schema"
)

type Parser struct {
	schema *schema.Schema
}

func (self *Parser) Parse(file map[string]any) error {
	err := self.domain(file)
	if err != nil {
		return err
	}

	err = self.metadata(file)
	if err != nil {
		return err
	}

	err = self.enums(file)
	if err != nil {
		return err
	}

	return nil
}

func (self *Parser) GetSchema() *schema.Schema {
	return self.schema
}

func NewParser() *Parser {
	return &Parser{
		schema: &schema.Schema{},
	}
}
