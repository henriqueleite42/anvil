package postgres

import (
	"fmt"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *hclFile) resolveSchema(
	entity *schemas.Entity,
) (string, error) {
	var dbSchema = "public"
	if entity.Schema != nil {
		dbSchema = *entity.Schema
	}

	if _, ok := self.dbSchemas[dbSchema]; !ok {
		self.dbSchemas[dbSchema] = fmt.Sprintf("schema \"%s\" {}", dbSchema)
	}

	return dbSchema, nil
}
