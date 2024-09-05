package postgres

import (
	"fmt"
	"strings"

	"github.com/henriqueleite42/anvil/cli/schemas"
)

type hclFile struct {
	schema    *schemas.Schema
	dbSchemas map[string]string
	enums     map[string]string
	tables    string
}

type SortedByOrder struct {
	Order int
	Key   string
}

func (self *hclFile) toString() string {
	schemasArr := []string{}
	for _, v := range self.dbSchemas {
		schemasArr = append(schemasArr, v)
	}
	schemas := strings.Join(schemasArr, "\n\n")

	enumsArr := []string{}
	for _, v := range self.enums {
		enumsArr = append(enumsArr, v)
	}
	enums := strings.Join(enumsArr, "\n\n")

	return fmt.Sprintf(`%s

%s

%s`, schemas, enums, self.tables)
}

func Parse(schema *schemas.Schema) (string, error) {
	proto := &hclFile{
		schema:    schema,
		dbSchemas: map[string]string{},
		enums:     map[string]string{},
	}

	err := proto.resolveTables(schema)
	if err != nil {
		return "", err
	}

	return proto.toString(), nil
}
