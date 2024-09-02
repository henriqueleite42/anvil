package parser

import (
	"fmt"
)

func (self *anvToAnvpParser) domain(file map[string]any) error {
	path := "Domain"

	domainAny, ok := file["Domain"]
	if !ok {
		return fmt.Errorf("\"%s\" must be specified", path)
	}
	domainString, ok := domainAny.(string)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `string`", path)
	}

	self.schema.Domain = domainString

	return nil
}
