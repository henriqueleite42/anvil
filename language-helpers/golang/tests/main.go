package main

import (
	"log"

	formatter_test "github.com/henriqueleite42/anvil/language-helpers/golang/tests/formatter"
	types_test "github.com/henriqueleite42/anvil/language-helpers/golang/tests/types"
)

func main() {
	formatter_test.PascalToSnake()
	formatter_test.PascalToKebab()
	formatter_test.PascalToCamel()

	types_test.ParseType()

	log.Default().Println("All done")
}
