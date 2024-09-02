package internal

import (
	"fmt"
	"strings"

	"github.com/henriqueleite42/anvil/cli/schemas"
)

type protoFile struct {
	schema   *schemas.Schema
	imports  map[string]bool
	enums    map[string]string
	messages map[string]string
	service  string
}

func (self *protoFile) toString() string {
	importsPaths := []string{}
	for k := range self.imports {
		importsPaths = append(importsPaths, k)
	}
	imports := strings.Join(importsPaths, "\n")

	enumsTypes := []string{}
	for _, v := range self.enums {
		enumsTypes = append(enumsTypes, v)
	}
	enums := strings.Join(enumsTypes, "\n\n")

	messagesTypes := []string{}
	for _, v := range self.messages {
		messagesTypes = append(messagesTypes, v)
	}
	messages := strings.Join(messagesTypes, "\n\n")

	return fmt.Sprintf(`syntax = "proto3";

%s

%s

%s

%s

`, imports, self.service, enums, messages)
}

func Parse(schema *schemas.Schema) (string, error) {
	proto := &protoFile{
		schema:   schema,
		imports:  map[string]bool{},
		enums:    map[string]string{},
		messages: map[string]string{},
	}

	err := proto.resolveService(schema)
	if err != nil {
		return "", err
	}

	return proto.toString(), nil
}
