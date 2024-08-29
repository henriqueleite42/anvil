package parser

import (
	"fmt"
	"strings"

	"github.com/anvil/anvil/internal/hashing"
	"github.com/anvil/anvil/internal/schema"
)

func (self *anvToAnvpParser) resolveRelationship(path string, k string, v any) (string, error) {
	if self.schema.Relationships == nil {
		self.schema.Relationships = &schema.Relationships{}
	}
	if self.schema.Relationships.Relationships == nil {
		self.schema.Relationships.Relationships = map[string]*schema.Relationship{}
	}

	originalPath := path + "." + k
	originalPathHash := hashing.String(originalPath)

	_, ok := self.schema.Relationships.Relationships[originalPathHash]
	if ok {
		return originalPathHash, nil
	}

	vMap, ok := v.(map[string]any)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s\" to `map[string]any`", path, k)
	}

	uriAny, ok := vMap["Uri"]
	if !ok {
		return "", fmt.Errorf("\"Uri\" is a required property to \"%s.%s\"", path, k)
	}
	uriString, ok := uriAny.(string)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s.Uri\" to `string`", path, k)
	}

	// Same project domains should be in the same repository
	var isSameProject bool
	if strings.HasPrefix(uriString, "./") || strings.HasPrefix(uriString, "../") {
		isSameProject = true
	}

	rootNode, err := getRootNode(originalPath)
	if err != nil {
		return "", err
	}

	relationship := &schema.Relationship{
		Name:          k,
		RootNode:      rootNode,
		OriginalPath:  originalPath,
		Uri:           uriString,
		Version:       "",
		IsSameProject: isSameProject,
	}

	stateHash, err := hashing.Struct(relationship)
	if err != nil {
		return "", fmt.Errorf("fail to get relationship \"%s\" state hash", originalPath)
	}

	relationship.StateHash = stateHash
	self.schema.Relationships.Relationships[originalPathHash] = relationship

	return originalPathHash, nil
}

func (self *anvToAnvpParser) relationships(file map[string]any) error {
	relationshipsSchema, ok := file["Relationships"]
	if !ok {
		return nil
	}

	fullPath := self.getPath("Relationships")

	relationshipsMap, ok := relationshipsSchema.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `map[string]any`", fullPath)
	}

	for k, v := range relationshipsMap {
		relationshipId, err := self.resolveRelationship(fullPath, k, v)
		if err != nil {
			return err
		}

		// Load relationships into the same schema
		relationship := self.schema.Relationships.Relationships[relationshipId]

		formattedUri := relationship.Uri
		if strings.HasPrefix(formattedUri, ".") {
			path := strings.Split(self.filePath, "/")
			pathWithoutFile := path[0 : len(path)-1]
			formattedUri = strings.Join(pathWithoutFile, "/") + "/" + relationship.Uri
		}

		relationShipFile, err := self.readAnvFile(formattedUri)
		if err != nil {
			return err
		}

		parser := &anvToAnvpParser{
			schema:   self.schema,
			basePath: relationship.OriginalPath,
		}

		err = parser.parse(relationShipFile)
		if err != nil {
			return err
		}
	}

	return nil
}
