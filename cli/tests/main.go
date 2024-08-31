package main

import parser_anv_test "github.com/anvil/anvil/tests/parser_anv"

func main() {
	logJson := false

	parser_anv_test.Authentication(true)
	parser_anv_test.EmailMailer(logJson)
	parser_anv_test.UrlShortener(logJson)
	parser_anv_test.Counter(logJson)
}
