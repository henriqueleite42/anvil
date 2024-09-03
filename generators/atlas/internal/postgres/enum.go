package postgres

import (
	"fmt"
	"strings"

	"github.com/henriqueleite42/anvil/cli/schemas"
)

func (self *hclFile) resolveEnum(
	schema *schemas.Schema,
	dbSchema string,
	enumHash *string,
) (string, error) {
	if enumHash == nil {
		return "", fmt.Errorf("missing EnumHash")
	}
	if schema.Enums == nil || schema.Enums.Enums == nil {
		return "", fmt.Errorf("no enums found in schema, but required for \"%s\"", *enumHash)
	}

	enum, ok := schema.Enums.Enums[*enumHash]
	if !ok {
		return "", fmt.Errorf("enums \"%s\" not found in schema", *enumHash)
	}

	enumValuesArr := []string{}
	for _, v := range enum.Values {
		enumValuesArr = append(enumValuesArr, fmt.Sprintf("		\"%s\"", v.Value))
	}
	enumValues := strings.Join(enumValuesArr, ",\n")

	self.enums[enum.Name] = fmt.Sprintf(`enum "%s" {
	schema = %s
	values = [
%s
	]
}`, enum.Name, dbSchema, enumValues)

	return enum.Name, nil
}
