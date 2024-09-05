package contract

import (
	"fmt"
	"sort"
	"strings"

	"github.com/henriqueleite42/anvil/cli/formatter"
	"github.com/henriqueleite42/anvil/cli/schemas"
)

type contractFile struct {
	schema  *schemas.Schema
	enums   map[string]string
	types   map[string]string
	methods []string
	imports map[string]bool
}

type SortedByOrder struct {
	Order int
	Key   string
}

func (self *contractFile) toString() string {
	nameSnake := formatter.PascalToSnake(self.schema.Domain)

	sortedImports := []string{}
	for k := range self.imports {
		sortedImports = append(sortedImports, k)
	}
	sort.Slice(sortedImports, func(i, j int) bool {
		return sortedImports[i] < sortedImports[j]
	})
	var imports string
	for _, v := range sortedImports {
		imports += "\n" + v
	}

	enumsTypes := []string{}
	for _, v := range self.enums {
		enumsTypes = append(enumsTypes, v)
	}
	enums := strings.Join(enumsTypes, "\n")

	typesTypes := []string{}
	for _, v := range self.types {
		typesTypes = append(typesTypes, v)
	}
	types := strings.Join(typesTypes, "\n")

	methods := strings.Join(self.methods, "\n")

	return fmt.Sprintf(`package %s

import (
%s
)

type ApiInput struct {
	Addr    string
	Timeout time.Duration
}

%s

%s

type %sApi interface {
	Close() error

%s
}
`, nameSnake, imports, enums, types, self.schema.Domain, methods)
}

func Parse(schema *schemas.Schema) (string, error) {
	contract := &contractFile{
		schema:  schema,
		enums:   map[string]string{},
		types:   map[string]string{},
		methods: []string{},
		imports: map[string]bool{
			"	\"time\"": true,
		},
	}

	err := contract.parseApi()
	if err != nil {
		return "", err
	}

	return contract.toString(), nil
}
