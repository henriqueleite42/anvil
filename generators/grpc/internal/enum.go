package internal

import (
	"strings"

	"github.com/henriqueleite42/anvil/generators/grpc/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *parser) resolveEnum(e *schemas.Enum) *templates.ProtofileTemplInputEnum {
	if existentEnum, ok := self.enums[e.Name]; ok {
		return existentEnum
	}

	result := &templates.ProtofileTemplInputEnum{
		Name:   e.Name,
		Values: make([]*templates.ProtofileTemplInputEnumValue, 0, len(e.Values)),
	}

	biggest := 0
	for _, v := range e.Values {
		newLen := len(v.Name)
		if newLen > biggest {
			biggest = newLen
		}
	}

	for k, v := range e.Values {
		result.Values = append(result.Values, &templates.ProtofileTemplInputEnumValue{
			Name:    v.Name,
			Spacing: strings.Repeat(" ", biggest-len(v.Name)),
			Idx:     k,
		})
	}

	self.enums[e.Name] = result

	return result
}
