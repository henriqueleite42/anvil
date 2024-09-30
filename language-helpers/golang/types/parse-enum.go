package types_parser

import (
	"strings"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *typeParser) ParseEnum(e *schemas.Enum) (*Enum, error) {
	if existentEnum, ok := self.enumsToAvoidDuplication[e.Name]; ok {
		return existentEnum, nil
	}

	var eType string
	if e.Type == schemas.EnumType_String {
		eType = "string"
	} else {
		eType = "int"
	}

	enum := &Enum{
		GolangName: e.Name,
		GolangType: eType,
		Values:     []*EnumValue{},
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

		enum.Values = append(enum.Values, &EnumValue{
			Idx:     k,
			Name:    v.Name,
			Spacing: strings.Repeat(" ", targetLen),
			Value:   v.Value,
		})
	}

	self.enumsToAvoidDuplication[e.Name] = enum
	self.enums = append(self.enums, enum)

	return enum, nil
}
