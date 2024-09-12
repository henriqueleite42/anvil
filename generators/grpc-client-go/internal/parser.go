package internal

import (
	"github.com/henriqueleite42/anvil/cli/schemas"
	"github.com/henriqueleite42/anvil/generators/grpc-client-go/internal/parser"
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
