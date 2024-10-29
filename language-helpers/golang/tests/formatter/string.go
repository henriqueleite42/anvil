package formatter_test

import (
	"log"

	"github.com/henriqueleite42/anvil/language-helpers/golang/formatter"
)

var pascalToSnakeTests = map[string]string{
	"FooBar":             "foo_bar",
	"FooBarBar":          "foo_bar_bar",
	"Foo123":             "foo123",
	"Foo123BarBar":       "foo123_bar_bar",
	"Foo123Bar123Bar":    "foo123_bar123_bar",
	"Foo123Bar123Bar123": "foo123_bar123_bar123",
	"Foo#Bar":            "foo#bar", // Wrong?
}

func PascalToSnake() {
	for k, v := range pascalToSnakeTests {
		r := formatter.PascalToSnake(k)
		if r != v {
			log.Fatalf("PascalToSnake: fail to parse Foobar.\nExpected: %s\nReceived: %s\n", v, r)
		}
	}
}

var pascalToKebabTests = map[string]string{
	"FooBar":             "foo-bar",
	"FooBarBar":          "foo-bar-bar",
	"Foo123":             "foo123",
	"Foo123BarBar":       "foo123-bar-bar",
	"Foo123Bar123Bar":    "foo123-bar123-bar",
	"Foo123Bar123Bar123": "foo123-bar123-bar123",
	"Foo#Bar":            "foo#bar", // Wrong?
}

func PascalToKebab() {
	for k, v := range pascalToKebabTests {
		r := formatter.PascalToKebab(k)
		if r != v {
			log.Fatalf("PascalToKebab: fail to parse Foobar.\nExpected: %s\nReceived: %s\n", v, r)
		}
	}
}

var pascalToCamelTests = map[string]string{
	"FooBar":             "fooBar",
	"FooBarBar":          "fooBarBar",
	"Foo123":             "foo123",
	"Foo123BarBar":       "foo123BarBar",
	"Foo123Bar123Bar":    "foo123Bar123Bar",
	"Foo123Bar123Bar123": "foo123Bar123Bar123",
	"Foo#Bar":            "foo#Bar",
}

func PascalToCamel() {
	for k, v := range pascalToCamelTests {
		r := formatter.PascalToCamel(k)
		if r != v {
			log.Fatalf("PascalToCamel: fail to parse Foobar.\nExpected: %s\nReceived: %s\n", v, r)
		}
	}
}
