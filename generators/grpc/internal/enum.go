package internal

import (
	"fmt"
	"strings"

	"github.com/henriqueleite42/anvil/cli/schemas"
)

func (self *protoFile) resolveEnum(e *schemas.Enum) string {
	if _, ok := self.enums[e.Name]; ok {
		return e.Name
	}

	values := []string{}
	for k, v := range e.Values {
		values = append(values, fmt.Sprintf("	%s = %d;", v.Value, k))
	}

	self.enums[e.Name] = fmt.Sprintf(`enum %s {
%s
}`, e.Name, strings.Join(values, "\n"))

	return e.Name
}
