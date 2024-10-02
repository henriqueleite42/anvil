package parser

import (
	"fmt"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *Parser) ResolveEntity(e *schemas.Entity) error {
	if e == nil {
		return fmt.Errorf("entity must be specified, received nil")
	}

	if self.Schema == nil {
		return fmt.Errorf("missing schema")
	}
	if self.Schema.Types == nil || self.Schema.Types.Types == nil {
		return fmt.Errorf("missing schema types")
	}

	if e.Columns == nil {
		return fmt.Errorf("entity \"%s\" missing required field \"Columns\"", e.Name)
	}
	lenColumns := len(e.Columns)
	if lenColumns == 0 {
		return fmt.Errorf("entity \"%s\" has to have at least 1 \"Columns\"", e.Name)
	}

	t, ok := self.Schema.Types.Types[e.TypeHash]
	if !ok {
		return fmt.Errorf("type \"%s\" for entity \"%s\" not found", e.TypeHash, e.Name)
	}

	self.GoTypesParserModels.ParseType(t, nil)

	return nil
}
