package generator_grpc

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/anvil/anvil/cmd/schema"
)

const protofileTempl = `
syntax = "proto3";

{{imports}}

{{options}}

service {{domain}} {
{{rpcs}}
}

{{messages}}
`

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
var emptyOutputImport = `import "google/protobuf/empty.proto";`

func toKebabCase(str string) string {
	kebab := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	kebab = matchAllCap.ReplaceAllString(kebab, "${1}-${2}")
	return strings.ToLower(kebab)
}

func replacePlaceholders(text string, placeholders map[string]string) (string, error) {
	replacedText := text

	for k, v := range placeholders {
		re, err := regexp.Compile("{{" + k + "}}")
		if err != nil {
			return "", errors.New("fail to compile regex of placeholder")
		}

		replacedText = re.ReplaceAllString(replacedText, v)
	}

	return replacedText, nil
}

func Generate(schema *schema.Schema) error {
	importsMap := map[string]bool{}
	optionsSlice := []string{
		// Temporary, only because we only generate gRPC files for Golang projects in a specific pattern, MVP
		`option go_package = "./internal/delivery/grpc/pb";`,
	}
	rpcsSlice := []string{}
	messagesSlice := []string{}

	for k, v := range schema.Usecase.Methods {
		if v.Delivery == nil || v.Delivery.Grpc == nil {
			continue
		}

		input := ""
		output := "google.protobuf.Empty"

		if v.Input != nil {
			input = k + "Input"
		}

		if v.Output != nil {
			output = k + "Output"
		} else {
			if _, ok := importsMap[emptyOutputImport]; !ok {
				importsMap[emptyOutputImport] = true
			}
		}

		rpcsSlice = append(rpcsSlice, fmt.Sprintf("	rpc %s(%s) returns (%s);", k, input, output))
	}

	if len(rpcsSlice) == 0 {
		return nil
	}

	importsSlice := []string{}

	for k, _ := range importsMap {
		importsSlice = append(importsSlice, k)
	}

	options := strings.Join(optionsSlice, "\n")
	rpcs := strings.Join(rpcsSlice, "\n")
	imports := strings.Join(importsSlice, "\n")
	messages := strings.Join(messagesSlice, "\n")

	finalValue, err := replacePlaceholders(protofileTempl, map[string]string{
		"imports":  imports,
		"options":  options,
		"domain":   schema.Domain,
		"rpcs":     rpcs,
		"messages": messages,
	})
	if err != nil {
		return err
	}

	kebabDomain := toKebabCase(schema.Domain)

	path := "./pkg/" + kebabDomain + ".proto"
	os.MkdirAll(filepath.Dir(path), 0700)
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return err
	}

	file.Write([]byte(finalValue))

	return nil
}
