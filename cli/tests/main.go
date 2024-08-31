package main

import (
	files_test "github.com/henriqueleite42/anvil/cli/tests/files"
	parser_anv_test "github.com/henriqueleite42/anvil/cli/tests/parser_anv"
)

func main() {
	logJson := false

	files_test.ReadAnvpFile(logJson)

	parser_anv_test.Authentication(logJson)
	parser_anv_test.EmailMailer(logJson)
	parser_anv_test.UrlShortener(logJson)
	parser_anv_test.Counter(logJson)
}
