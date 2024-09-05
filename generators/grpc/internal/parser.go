package internal

import (
	"fmt"
	"sort"

	"github.com/henriqueleite42/anvil/cli/schemas"
)

type protoFile struct {
	schema   *schemas.Schema
	imports  map[string]bool
	enums    map[string]string
	messages map[string]string
	service  string
}

type SortedByOrder struct {
	Order int
	Key   string
}

func (self *protoFile) toString() string {
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

	service := self.service
	if imports != "" {
		service = "\n" + service
	}

	sortedEnums := []string{}
	for _, v := range self.enums {
		sortedEnums = append(sortedEnums, v)
	}
	sort.Slice(sortedEnums, func(i, j int) bool {
		return sortedEnums[i] < sortedEnums[j]
	})
	var enums string
	for _, v := range sortedEnums {
		enums += "\n" + v
	}

	sortedMessages := []string{}
	for _, v := range self.messages {
		sortedMessages = append(sortedMessages, v)
	}
	sort.Slice(sortedMessages, func(i, j int) bool {
		return sortedMessages[i] < sortedMessages[j]
	})
	var messages string
	for _, v := range sortedMessages {
		messages += "\n" + v
	}

	if enums != "" {
		messages = "\n" + messages
	}

	return fmt.Sprintf(`syntax = "proto3";
%s
%s
%s%s`, imports, service, enums, messages)
}

func Parse(schema *schemas.Schema) (string, error) {
	proto := &protoFile{
		schema:   schema,
		imports:  map[string]bool{},
		enums:    map[string]string{},
		messages: map[string]string{},
	}

	err := proto.resolveService()
	if err != nil {
		return "", err
	}

	return proto.toString(), nil
}
