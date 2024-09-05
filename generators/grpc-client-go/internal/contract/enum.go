package contract

import (
	"fmt"
	"strings"

	"github.com/henriqueleite42/anvil/cli/schemas"
	"github.com/henriqueleite42/anvil/generators/grpc-client-go/internal/spacer"
)

func (self *contractFile) resolveEnum(e *schemas.Enum) (string, error) {
	if _, ok := self.enums[e.Name]; ok {
		return e.Name, nil
	}

	values, err := spacer.Space(
		e.Values,
		func(v *schemas.EnumValue) ([]string, error) {
			var value string
			if e.Type == "String" {
				value = fmt.Sprintf("\"%s\"", v.Value)
			} else {
				value = v.Value
			}

			return []string{
				e.Name + "_" + v.Name,
				value,
			}, nil
		},
		func(s []string, i int) (string, error) {
			targetLen := i - len(s[0])
			return fmt.Sprintf(
				"	%s %s = %s",
				s[0]+strings.Repeat(" ", targetLen),
				e.Name,
				s[1],
			), nil
		},
	)
	if err != nil {
		return "", err
	}

	var typeEnum string
	if e.Type == "String" {
		typeEnum = "string"
	} else {
		typeEnum = "int"
	}

	self.enums[e.Name] = &ItemWithOrder{
		Order: len(self.enums),
		Value: fmt.Sprintf(`type %s %s

const (
%s
)`, e.Name, typeEnum, strings.Join(values, "\n")),
	}

	return e.Name, nil
}
