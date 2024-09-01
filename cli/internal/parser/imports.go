package parser

import (
	"fmt"

	"github.com/henriqueleite42/anvil/cli/internal/hashing"
	"github.com/henriqueleite42/anvil/cli/schemas"
)

func (self *anvToAnvpParser) resolveImport(i *resolveInput) (string, error) {
	if self.schema.Imports == nil {
		self.schema.Imports = &schemas.Imports{}
	}
	if self.schema.Imports.Imports == nil {
		self.schema.Imports.Imports = map[string]*schemas.Import{}
	}

	ref := self.getRef(i.ref, "Imports."+i.k)
	refHash := hashing.String(ref)

	_, ok := self.schema.Imports.Imports[refHash]
	if ok {
		return refHash, nil
	}

	vMap, ok := i.v.(map[string]any)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s\" to `map[string]any`", i.path, i.k)
	}

	var importImport *schemas.ImportImport = nil
	importAny, ok := vMap["Import"]
	if ok {
		importMap, ok := importAny.(map[string]any)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Import\" to `map[string]any`", i.path, i.k)
		}

		var aliasPointer *string = nil
		aliasAny, ok := importMap["Alias"]
		if ok {
			aliasString, ok := aliasAny.(string)
			if !ok {
				return "", fmt.Errorf("fail to parse \"%s.%s.Import.Alias\" to `string`", i.path, i.k)
			}
			aliasPointer = &aliasString
		}

		pathAny, ok := importMap["Path"]
		if !ok {
			return "", fmt.Errorf("\"Path\" is a required property to \"%s.%s.Import\"", i.path, i.k)
		}
		pathString, ok := pathAny.(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Import.Path\" to `string`", i.path, i.k)
		}

		importImport = &schemas.ImportImport{
			Alias: aliasPointer,
			Path:  pathString,
		}
	}

	typeAny, ok := vMap["Type"]
	if !ok {
		return "", fmt.Errorf("\"Type\" is a required property to \"%s.%s\"", i.path, i.k)
	}
	typeString, ok := typeAny.(string)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s.Type\" to `string`", i.path, i.k)
	}

	rootNode, err := getRootNode(i.path)
	if err != nil {
		return "", err
	}

	schemaImports := &schemas.Import{
		OriginalPath: fmt.Sprintf("%s.%s", i.path, i.k),
		Name:         i.k,
		RootNode:     rootNode,
		Import:       importImport,
		Type:         typeString,
	}

	stateHash, err := hashing.Struct(schemaImports)
	if err != nil {
		return "", fmt.Errorf("fail to get import \"%s.%s\" state hash", i.path, i.k)
	}

	schemaImports.StateHash = stateHash
	self.schema.Imports.Imports[refHash] = schemaImports

	return refHash, nil
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
		_, err := self.resolveImport(&resolveInput{
			path: fullPath,
			ref:  "",
			k:    k,
			v:    v,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
