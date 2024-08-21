package parser

import (
	"errors"
)

func (self *Parser) domain(file map[string]any) error {
	domainSchema, ok := file["Domain"]
	if !ok {
		return errors.New("\"Domain\" must be specified")
	}

	domainString, ok := domainSchema.(string)
	if !ok {
		return errors.New("fail to parse \"Domain\" to `string`")
	}

	self.schema.Domain = domainString

	return nil
}
