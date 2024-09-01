package main

import (
	files_test "github.com/henriqueleite42/anvil/cli/tests/files"
	parser_test "github.com/henriqueleite42/anvil/cli/tests/parser"
)

func main() {
	logJson := false

	files_test.ReadAnvpFile(logJson)

	parser_test.Authentication(logJson)
	parser_test.EmailMailer(logJson)
	parser_test.UrlShortener(logJson)
	parser_test.Counter(logJson)
}
