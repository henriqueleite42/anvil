package parser

import (
	"fmt"
)

func (self *anvToAnvpParser) domain(file map[string]any) error {
	domainSchema, ok := file["Domain"]
	if !ok {
		return fmt.Errorf("\"%s\" must be specified", self.getPath("Domain"))
	}

	fullPath := self.getPath("Domain")

	domainString, ok := domainSchema.(string)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `string`", fullPath)
	}

	self.schema.Domain = domainString

	return nil
}
