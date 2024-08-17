package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/anvil/anvil/internal/hashing"
	"github.com/anvil/anvil/internal/schema"
)

func (self *Parser) resolveRelationship(path string, k string, v any) (string, error) {
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

	originalPath := path + "." + k
	originalPathHash := hashing.String(originalPath)

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

	if self.schema.Relationships == nil {
		self.schema.Relationships = &schema.Relationships{}
	}
	if self.schema.Relationships.Relationships == nil {
		self.schema.Relationships.Relationships = map[string]*schema.Relationship{}
	}

	relationship.StateHash = stateHash
	self.schema.Relationships.Relationships[originalPathHash] = relationship

	return originalPathHash, nil
}

func (self *Parser) relationships(file map[string]any) error {
	relationshipsSchema, ok := file["Relationships"]
	if !ok {
		return nil
	}

	relationshipsMap, ok := relationshipsSchema.(map[string]any)
	if !ok {
		return errors.New("fail to parse \"Relationships\" to `map[string]any`")
	}

	for k, v := range relationshipsMap {
		_, err := self.resolveRelationship("Relationships", k, v)
		if err != nil {
			return err
		}
	}

	return nil
}
