package main

import parser_test "github.com/anvil/anvil/tests/parser"

func main() {
	logJson := false

	parser_test.Authentication(true)
	parser_test.EmailMailer(logJson)
	parser_test.UrlShortener(logJson)
	parser_test.Counter(logJson)
}
