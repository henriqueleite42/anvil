package parser

import (
	"fmt"

	"github.com/anvil/anvil/internal/hashing"
	"github.com/anvil/anvil/internal/schema"
)

func (self *anvToAnvpParser) resolveImport(path string, k string, v any) (string, error) {
	if self.schema.Imports == nil {
		self.schema.Imports = &schema.Imports{}
	}
	if self.schema.Imports.Imports == nil {
		self.schema.Imports.Imports = map[string]*schema.Import{}
	}

	originalPath := path + "." + k
	originalPathHash := hashing.String(originalPath)

	_, ok := self.schema.Imports.Imports[originalPathHash]
	if ok {
		return originalPathHash, nil
	}

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

	rootNode, err := getRootNode(path)
	if err != nil {
		return "", err
	}

	schemaImports := &schema.Import{
		Name:         k,
		RootNode:     rootNode,
		OriginalPath: originalPath,
		Import:       importImport,
		Type:         typeString,
	}

	stateHash, err := hashing.Struct(schemaImports)
	if err != nil {
		return "", fmt.Errorf("fail to get import \"%s\" state hash", originalPath)
	}

	schemaImports.StateHash = stateHash
	self.schema.Imports.Imports[originalPathHash] = schemaImports

	return originalPathHash, nil
}

func (self *anvToAnvpParser) imports(file map[string]any) error {
	importsSchema, ok := file["Imports"]
	if !ok {
		return nil
	}

	fullPath := self.getPath("Imports")

	importsMap, ok := importsSchema.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `map[string]any`", fullPath)
	}

	for k, v := range importsMap {
		_, err := self.resolveImport(fullPath, k, v)
		if err != nil {
			return err
		}
	}

	return nil
}
