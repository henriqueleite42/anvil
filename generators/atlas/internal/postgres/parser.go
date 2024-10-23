package postgres

import (
	"fmt"

	"github.com/henriqueleite42/anvil/generators/atlas/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	"github.com/henriqueleite42/anvil/language-helpers/golang/template"
)

func Parse(schema *schemas.AnvpSchema, silent bool, outputFolderPath *string) error {
	if schema.Entities == nil || schema.Entities.Entities == nil {
		return fmt.Errorf("no entities to create tables")
	}

	templateManager := template.NewTemplateManager()
	err := templateManager.AddTemplate("hcl", templates.HclTempl)
	if err != nil {
		return err
	}

	enums, err := resolveEnums(schema)
	if err != nil {
		return err
	}

	entities, err := resolveEntities(schema)
	if err != nil {
		return err
	}

	templInput := &templates.HclTemplInput{
		Enums:    enums,
		Entities: entities,
	}

	hclFile, err := templateManager.Parse("hcl", templInput)
	if err != nil {
		return err
	}

	err = WriteFile(schema, outputFolderPath, hclFile)
	if err != nil {
		return err
	}

	return nil
}
