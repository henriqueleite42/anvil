package parser_anv

import (
	"fmt"
)

func (self *anvToAnvpParser) domain(file map[string]any) error {
	path := self.getPath("Domain")

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
