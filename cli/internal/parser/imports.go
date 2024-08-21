package parser

import (
	"errors"
	"fmt"

	"github.com/anvil/anvil/internal/hashing"
	"github.com/anvil/anvil/internal/schema"
)

func (self *Parser) resolveImport(path string, k string, v any) (string, error) {
	vMap, ok := v.(map[string]any)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s\" to `map[string]any`", path, k)
	}

	var importImport *schema.ImportImport = nil
	importAny, ok := vMap["Import"]
	if ok {
		importMap, ok := importAny.(map[string]any)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Import\" to `map[string]any`", path, k)
		}

		var aliasPointer *string = nil
		aliasAny, ok := importMap["Alias"]
		if ok {
			aliasString, ok := aliasAny.(string)
			if !ok {
				return "", fmt.Errorf("fail to parse \"%s.%s.Import.Alias\" to `string`", path, k)
			}
			aliasPointer = &aliasString
		}

		pathAny, ok := importMap["Path"]
		if !ok {
			return "", fmt.Errorf("\"Path\" is a required property to \"%s.%s.Import\"", path, k)
		}
		pathString, ok := pathAny.(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Import.Path\" to `string`", path, k)
		}

		importImport = &schema.ImportImport{
			Alias: aliasPointer,
			Path:  pathString,
		}
	}

	typeAny, ok := vMap["Type"]
	if !ok {
		return "", fmt.Errorf("\"Type\" is a required property to \"%s.%s\"", path, k)
	}
	typeString, ok := typeAny.(string)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s.Type\" to `string`", path, k)
	}

	originalPath := path + "." + k
	originalPathHash := hashing.String(originalPath)

	rootNode, err := getRootNode(path)
	if err != nil {
		return "", err
	}

	import_ := &schema.Import{
		Name:         k,
		RootNode:     rootNode,
		OriginalPath: originalPath,
		Import:       importImport,
		Type:         typeString,
	}

	stateHash, err := hashing.Struct(import_)
	if err != nil {
		return "", fmt.Errorf("fail to get import \"%s\" state hash", originalPath)
	}

	if self.schema.Relationships == nil {
		self.schema.Relationships = &schema.Relationships{}
	}
	if self.schema.Relationships.Relationships == nil {
		self.schema.Relationships.Relationships = map[string]*schema.Relationship{}
	}

	import_.StateHash = stateHash
	self.schema.Imports.Imports[originalPathHash] = import_

	return originalPathHash, nil
}

func (self *Parser) imports(file map[string]any) error {
	importsSchema, ok := file["Imports"]
	if !ok {
		return nil
	}

	importsMap, ok := importsSchema.(map[string]any)
	if !ok {
		return errors.New("fail to parse \"Imports\" to `map[string]any`")
	}

	for k, v := range importsMap {
		_, err := self.resolveImport("Imports", k, v)
		if err != nil {
			return err
		}
	}

	return nil
}
