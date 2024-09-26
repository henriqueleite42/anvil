package parser

import (
	"strings"

	"github.com/henriqueleite42/anvil/generators/grpc-client-go/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *parserManager) toEnum(e *schemas.Enum) (*templates.TemplEnum, error) {
	if existentEnum, ok := self.enums[e.Name]; ok {
		return existentEnum, nil
	}

	var eType string
	if e.Type == schemas.EnumType_String {
		eType = "string"
	} else {
		eType = "int"
	}

	enum := &templates.TemplEnum{
		Name:   e.Name,
		Type:   eType,
		Values: []*templates.TemplEnumValue{},
	}

	biggest := len(e.Values[0].Name)
	for _, v := range e.Values {
		newLen := len(v.Name)
		if newLen > biggest {
			biggest = newLen
		}
	}

	for k, v := range e.Values {
		targetLen := biggest - len(v.Name)

		enum.Values = append(enum.Values, &templates.TemplEnumValue{
			Idx:     k,
			Name:    v.Name,
			Spacing: strings.Repeat(" ", targetLen),
			Value:   v.Value,
		})
	}

	self.enums[e.Name] = enum

	return enum, nil
}
