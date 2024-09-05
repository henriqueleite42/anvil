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
	enums   map[string]*ItemWithOrder
	types   map[string]*ItemWithOrder
	methods []string
	imports map[string]bool
}

type SortedByOrder struct {
	Order int
	Key   string
}

type ItemWithOrder struct {
	Value string
	Order int
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

	sortedEnums := []*SortedByOrder{}
	for k, v := range self.enums {
		sortedEnums = append(sortedEnums, &SortedByOrder{
			Order: v.Order,
			Key:   k,
		})
	}
	sort.Slice(sortedEnums, func(i, j int) bool {
		return sortedEnums[i].Order < sortedEnums[j].Order
	})
	enumsTypes := []string{}
	for _, v := range sortedEnums {
		enumsTypes = append(enumsTypes, self.enums[v.Key].Value)
	}
	enums := strings.Join(enumsTypes, "\n")

	sortedTypes := []*SortedByOrder{}
	for k, v := range self.types {
		sortedTypes = append(sortedTypes, &SortedByOrder{
			Order: v.Order,
			Key:   k,
		})
	}
	sort.Slice(sortedTypes, func(i, j int) bool {
		return sortedTypes[i].Order < sortedTypes[j].Order
	})
	typesTypes := []string{}
	for _, v := range sortedTypes {
		typesTypes = append(typesTypes, self.types[v.Key].Value)
	}
	types := strings.Join(typesTypes, "\n")

	methods := strings.Join(self.methods, "\n")

	return fmt.Sprintf(`package %s

import (%s
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
		enums:   map[string]*ItemWithOrder{},
		types:   map[string]*ItemWithOrder{},
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
