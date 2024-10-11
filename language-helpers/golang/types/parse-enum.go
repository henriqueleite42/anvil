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
		eType = "int32"
	}

	enum := &Enum{
		GolangPkg:        self.enumsPkg,
		GolangName:       e.Name,
		GolangType:       eType,
		Values:           make([]*EnumValue, 0, len(e.Values)),
		DeprecatedValues: []*EnumValue{},
	}

	biggest := len(e.Values[0].Name)
	for _, v := range e.Values {
		newLen := len(v.Name)
		if newLen > biggest {
			biggest = newLen
		}
	}

	for _, v := range e.Values {
		targetLen := biggest - len(v.Name)

		value := &EnumValue{
			Idx:     v.Index,
			Name:    v.Name,
			Spacing: strings.Repeat(" ", targetLen),
			Value:   v.Value,
		}

		if v.Deprecated {
			enum.DeprecatedValues = append(enum.DeprecatedValues, value)
		} else {
			enum.Values = append(enum.Values, value)
		}
	}

	self.enumsToAvoidDuplication[e.Name] = enum
	self.enums = append(self.enums, enum)

	return enum, nil
}
