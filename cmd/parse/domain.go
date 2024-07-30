package parse

import (
	"errors"

	"github.com/anuntech/hephaestus/cmd/types"
)

func Domain(s *types.Schema, yaml map[string]any) error {
	domain, ok := yaml["Domain"]
	if !ok {
		return errors.New("Domain must be specified")
	}
	s.Domain = domain.(string)

	return nil
}
