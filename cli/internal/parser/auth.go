package parser

import (
	"fmt"

	"github.com/anvil/anvil/internal/hashing"
	"github.com/anvil/anvil/internal/schema"
)

func (self *anvToAnvpParser) auth(file map[string]any) error {
	authSchema, ok := file["Auth"]
	if !ok {
		return nil
	}

	fullPath := self.getPath("Auth")

	authMap, ok := authSchema.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `map[string]any`", fullPath)
	}

	if self.schema.Auths == nil {
		self.schema.Auths = &schema.Auths{}
	}
	if self.schema.Auths.Auths == nil {
		self.schema.Auths.Auths = map[string]*schema.Auth{}
	}

	for k, v := range authMap {
		vMap, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("fail to parse \"%s.%s\" to `map[string]any`", fullPath, k)
		}

		var description *string = nil
		descriptionAny, ok := vMap["Description"]
		if ok {
			descriptionString, ok := descriptionAny.(string)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.%s.Description\" to `string`", fullPath, k)
			}
			description = &descriptionString
		}

		var scheme string
		schemeAny, ok := vMap["Scheme"]
		if ok {
			schemeString, ok := schemeAny.(string)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.%s.Scheme\" to `string`", fullPath, k)
			}
			scheme = schemeString
		}

		var format *string = nil
		formatAny, ok := vMap["Format"]
		if ok {
			formatString, ok := formatAny.(string)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.%s.Format\" to `string`", fullPath, k)
			}
			format = &formatString
		}

		var applyToallRoutes bool
		applyToallRoutesAny, ok := vMap["ApplyToAllRoutes"]
		if ok {
			applyToallRoutesString, ok := applyToallRoutesAny.(bool)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.%s.ApplyToAllRoutes\" to `string`", fullPath, k)
			}
			applyToallRoutes = applyToallRoutesString
		}

		originalPath := fullPath + "." + k
		originalPathHash := hashing.String(originalPath)

		rootNode, err := getRootNode(fullPath)
		if err != nil {
			return err
		}

		auth := &schema.Auth{
			Name:             k,
			RootNode:         rootNode,
			OriginalPath:     originalPath,
			Description:      description,
			Scheme:           scheme,
			Format:           format,
			ApplyToAllRoutes: applyToallRoutes,
		}

		stateHash, err := hashing.Struct(auth)
		if err != nil {
			return fmt.Errorf("fail to get enum \"%s\" state hash", originalPath)
		}

		auth.StateHash = stateHash
		self.schema.Auths.Auths[originalPathHash] = auth
	}

	return nil
}
