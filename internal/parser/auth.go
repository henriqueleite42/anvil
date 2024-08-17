package parser

import (
	"errors"
	"fmt"

	"github.com/anvil/anvil/internal/hashing"
	"github.com/anvil/anvil/internal/schema"
)

func (self *Parser) auth(file map[string]any) error {
	authSchema, ok := file["Auth"]
	if !ok {
		return nil
	}

	authMap, ok := authSchema.(map[string]any)
	if !ok {
		return errors.New("fail to parse \"Auth\" to `map[string]any`")
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
			return fmt.Errorf("fail to parse \"Auth.%s\" to `map[string]any`", k)
		}

		var description *string = nil
		descriptionAny, ok := vMap["Description"]
		if ok {
			descriptionString, ok := descriptionAny.(string)
			if !ok {
				return fmt.Errorf("fail to parse \"Auth.%s.Description\" to `string`", k)
			}
			description = &descriptionString
		}

		var scheme string
		schemeAny, ok := vMap["Scheme"]
		if ok {
			schemeString, ok := schemeAny.(string)
			if !ok {
				return fmt.Errorf("fail to parse \"Auth.%s.Scheme\" to `string`", k)
			}
			scheme = schemeString
		}

		var format *string = nil
		formatAny, ok := vMap["Format"]
		if ok {
			formatString, ok := formatAny.(string)
			if !ok {
				return fmt.Errorf("fail to parse \"Auth.%s.Format\" to `string`", k)
			}
			format = &formatString
		}

		var applyToallRoutes bool
		applyToallRoutesAny, ok := vMap["ApplyToAllRoutes"]
		if ok {
			applyToallRoutesString, ok := applyToallRoutesAny.(bool)
			if !ok {
				return fmt.Errorf("fail to parse \"Auth.%s.ApplyToAllRoutes\" to `string`", k)
			}
			applyToallRoutes = applyToallRoutesString
		}

		originalPath := "Auth" + "." + k
		originalPathHash := hashing.String(originalPath)

		rootNode, err := getRootNode("Auth")
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
