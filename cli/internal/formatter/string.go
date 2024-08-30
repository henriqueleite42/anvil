package formatter

import (
	"regexp"
	"strings"
)

var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func PascalToSnake(str string) string {
	snake := matchAllCap.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}

func PascalToCamel(str string) string {
	return strings.ToLower(str[0:1]) + str[1:]
}
