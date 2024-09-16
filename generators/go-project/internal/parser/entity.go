package parser

import (
	"fmt"
	"sort"
	"strings"

	"github.com/henriqueleite42/anvil/cli/schemas"
	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
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

	result := &templates.TemplType{
		Name:         e.Name,
		OriginalType: t.Type,
		Props:        make([]*templates.TemplTypeProp, 0, lenColumns),
	}

	for _, v := range e.Columns {
		childTypeRef, ok := self.Schema.Types.Types[v.TypeHash]
		if !ok {
			return fmt.Errorf("child type \"%s\" of type \"%s\" not found", v.TypeHash, t.Name)
		}

		resolvedProp, err := self.ResolveMapProp(&ResolveMapPropInput{
			Kind:              Kind_Entity,
			Type:              childTypeRef,
			PrefixForChildren: result.Name,
			Tags: []string{
				fmt.Sprintf("db:\"%s\"", v.DbName),
			},
		})
		if err != nil {
			return err
		}

		result.Props = append(result.Props, resolvedProp)
	}
	sort.Slice(result.Props, func(i, j int) bool {
		return result.Props[i].Name < result.Props[j].Name
	})

	biggestPropName := 0
	biggestPropType := 0
	for _, v := range result.Props {
		if len(v.Name) > biggestPropName {
			biggestPropName = len(v.Name)
		}
		if len(v.Type) > biggestPropType {
			biggestPropType = len(v.Type)
		}
	}

	for _, v := range result.Props {
		v.Spacing1 = strings.Repeat(" ", biggestPropName-len(v.Name))
		v.Spacing2 = strings.Repeat(" ", biggestPropType-len(v.Type))
	}

	self.Entities = append(self.Entities, result)

	return nil
}
