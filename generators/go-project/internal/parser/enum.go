package parser

import (
	"fmt"
	"strings"

	"github.com/henriqueleite42/anvil/cli/schemas"
	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
)

func (self *Parser) resolveEnum(e *schemas.Enum) (*templates.TemplEnum, error) {
	if enumFound, ok := self.Enums[e.Name]; ok {
		return enumFound, nil
	}

	var enumType string
	if e.Type == schemas.EnumType_String {
		enumType = "string"
	} else if e.Type == schemas.EnumType_Int {
		enumType = "int32"
	} else {
		return nil, fmt.Errorf("fail to parse enum \"%s\": invalid type \"%s\"", e.Name, e.Type)
	}

	biggest := 0
	for _, v := range e.Values {
		newLen := len(v.Name)
		if newLen > biggest {
			biggest = newLen
		}
	}

	values := make([]*templates.TemplEnumValue, 0, len(e.Values))
	for k, v := range e.Values {
		values[k] = &templates.TemplEnumValue{
			Idx:     k,
			Name:    v.Name,
			Spacing: strings.Repeat(" ", biggest-len(v.Name)),
			Value:   v.Value,
		}
	}

	self.Enums[e.Name] = &templates.TemplEnum{
		Name:   e.Name,
		Type:   enumType,
		Values: values,
	}

	return self.Enums[e.Name], nil
}
