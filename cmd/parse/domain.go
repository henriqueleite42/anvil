package parse

import (
	"errors"

	"github.com/anvil/anvil/cmd/schema"
)

func domain(s *schema.Schema, yaml map[string]any) error {
	domain, ok := yaml["Domain"]
	if !ok {
		return errors.New("\"Domain\" must be specified")
	}
	s.Domain = domain.(string)

	return nil
}
