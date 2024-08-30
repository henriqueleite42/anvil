package parser

import (
	"fmt"
)

func (self *anvToAnvpParser) domain(file map[string]any) error {
	fullPath := self.getPath("Domain")

	domainSchema, ok := file["Domain"]
	if !ok {
		return fmt.Errorf("\"%s\" must be specified", fullPath)
	}

	domainString, ok := domainSchema.(string)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `string`", fullPath)
	}

	self.schema.Domain = domainString

	return nil
}
