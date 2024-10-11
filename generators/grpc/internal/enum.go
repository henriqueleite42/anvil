package internal

import (
	"fmt"
	"strings"

	"github.com/henriqueleite42/anvil/generators/grpc/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *parser) resolveEnum(e *schemas.Enum) *templates.ProtofileTemplInputEnum {
	if existentEnum, ok := self.enums[e.Name]; ok {
		return existentEnum
	}

	result := &templates.ProtofileTemplInputEnum{
		Name:             e.Name,
		Values:           make([]*templates.ProtofileTemplInputEnumValue, 0, len(e.Values)),
		DeprecatedValues: []int32{},
	}

	biggest := 0
	for _, v := range e.Values {
		name := fmt.Sprintf("%s_%s", e.Name, v.Name)
		newLen := len(name)
		if newLen > biggest {
			biggest = newLen
		}
	}

	for _, v := range e.Values {
		if v.Deprecated {
			result.DeprecatedValues = append(result.DeprecatedValues, v.Index)
			continue
		}

		name := fmt.Sprintf("%s_%s", e.Name, v.Name)
		result.Values = append(result.Values, &templates.ProtofileTemplInputEnumValue{
			Name:    name,
			Spacing: strings.Repeat(" ", biggest-len(name)),
			Idx:     v.Index,
		})
	}

	self.enums[e.Name] = result

	return result
}
