package parser_anv

import (
	"fmt"
	"strings"

	"github.com/anvil/anvil/internal/hashing"
	"github.com/anvil/anvil/schemas"
)

func (self *anvToAnvpParser) resolveRelationship(i *resolveInput) (string, error) {
	if self.schema.Relationships == nil {
		self.schema.Relationships = &schemas.Relationships{}
	}
	if self.schema.Relationships.Relationships == nil {
		self.schema.Relationships.Relationships = map[string]*schemas.Relationship{}
	}

	ref := self.getRef(i.ref, "Relationships."+i.k)
	refHash := hashing.String(ref)

	_, ok := self.schema.Relationships.Relationships[refHash]
	if ok {
		return refHash, nil
	}

	vMap, ok := i.v.(map[string]any)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s\" to `map[string]any`", i.path, i.k)
	}

	uriAny, ok := vMap["Uri"]
	if !ok {
		return "", fmt.Errorf("\"Uri\" is a required property to \"%s.%s\"", i.path, i.k)
	}
	uriString, ok := uriAny.(string)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s.Uri\" to `string`", i.path, i.k)
	}

	// Same project domains should be in the same repository
	var isSameProject bool
	if strings.HasPrefix(uriString, "./") || strings.HasPrefix(uriString, "../") {
		isSameProject = true
	}

	rootNode, err := getRootNode(i.path)
	if err != nil {
		return "", err
	}

	relationship := &schemas.Relationship{
		Ref:           ref,
		OriginalPath:  self.getPath(fmt.Sprintf("%s.%s", i.path, i.k)),
		Name:          i.k,
		RootNode:      rootNode,
		Uri:           uriString,
		Version:       "",
		IsSameProject: isSameProject,
	}

	stateHash, err := hashing.Struct(relationship)
	if err != nil {
		return "", fmt.Errorf("fail to get relationship \"%s.%s\" state hash", i.path, i.k)
	}

	relationship.StateHash = stateHash
	self.schema.Relationships.Relationships[refHash] = relationship

	return refHash, nil
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
		_, err := self.resolveRelationship(&resolveInput{
			path: fullPath,
			ref:  "",
			k:    k,
			v:    v,
		})
		if err != nil {
			return err
		}

		// TODO Load relationships
	}

	return nil
}
