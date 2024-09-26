package internal

import (
	"github.com/henriqueleite42/anvil/generators/grpc-client-go/internal/parser"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func Parse(schema *schemas.Schema, silent bool, outputFolderPath string) error {
	contract, implementation, err := parser.Parse(schema)
	if err != nil {
		return err
	}

	err = WriteFile(schema.Domain, outputFolderPath, "contract", contract)
	if err != nil {
		return err
	}

	err = WriteFile(schema.Domain, outputFolderPath, "implementation", implementation)
	if err != nil {
		return err
	}

	return nil
}
