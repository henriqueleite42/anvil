package types_parser

import (
	"fmt"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *typeParser) ParseEnum(e *schemas.Enum) (*Enum, error) {
	if existentEnum, ok := self.enumsToAvoidDuplication[e.Name]; ok {
		return existentEnum, nil
	}

	var eType string
	if e.Type == schemas.EnumType_String {
		eType = "string"
	}
	if e.Type == schemas.EnumType_Int {
		eType = "int"
	}
	if e.Type == schemas.EnumType_Int8 {
		eType = "int8"
	}
	if e.Type == schemas.EnumType_Int16 {
		eType = "int16"
	}
	if e.Type == schemas.EnumType_Int32 {
		eType = "int32"
	}
	if e.Type == schemas.EnumType_Int64 {
		eType = "int64"
	}
	if e.Type == schemas.EnumType_Uint {
		eType = "uint"
	}
	if e.Type == schemas.EnumType_Uint8 {
		eType = "uint8"
	}
	if e.Type == schemas.EnumType_Uint16 {
		eType = "uint16"
	}
	if e.Type == schemas.EnumType_Uint32 {
		eType = "uint32"
	}
	if e.Type == schemas.EnumType_Uint64 {
		eType = "uint64"
	}
	if eType == "" {
		return nil, fmt.Errorf("unsupported enum type: \"%s\"", e.Type)
	}

	enum := &Enum{
		AnvilEnum: e,

		Import:           self.getEnumsImport(e),
		GolangName:       e.Name,
		GolangType:       eType,
		Values:           make([]*EnumValue, 0, len(e.Values)),
		DeprecatedValues: []*EnumValue{},
	}

	for _, v := range e.Values {
		value := &EnumValue{
			Idx:   v.Index,
			Name:  v.Name,
			Value: v.Value,
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
