package formatter

import (
	"regexp"
	"strings"
	"unicode"
)

var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func PascalToSnake(str string) string {
	snake := matchAllCap.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}

func PascalToKebab(str string) string {
	snake := matchAllCap.ReplaceAllString(str, "${1}-${2}")
	return strings.ToLower(snake)
}

func PascalToCamel(str string) string {
	return strings.ToLower(str[0:1]) + str[1:]
}

func KebabToPascal(str string) string {
	var result strings.Builder

	capitalizeNext := true
	for _, r := range str {
		if r == '-' {
			capitalizeNext = true
		} else {
			if capitalizeNext {
				result.WriteRune(unicode.ToUpper(r))
				capitalizeNext = false
			} else {
				result.WriteRune(r)
			}
		}
	}

	return result.String()
}
