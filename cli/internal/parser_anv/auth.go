package parser_anv

import (
	"fmt"

	"github.com/henriqueleite42/anvil/cli/internal/hashing"
	"github.com/henriqueleite42/anvil/cli/schemas"
)

func (self *anvToAnvpParser) auth(file map[string]any) error {
	authSchema, ok := file["Auth"]
	if !ok {
		return nil
	}

	path := self.getPath("Auth")

	authMap, ok := authSchema.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `map[string]any`", path)
	}

	if self.schema.Auths == nil {
		self.schema.Auths = &schemas.Auths{}
	}
	if self.schema.Auths.Auths == nil {
		self.schema.Auths.Auths = map[string]*schemas.Auth{}
	}

	for k, v := range authMap {
		vMap, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("fail to parse \"%s.%s\" to `map[string]any`", path, k)
		}

		var description *string = nil
		descriptionAny, ok := vMap["Description"]
		if ok {
			descriptionString, ok := descriptionAny.(string)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.%s.Description\" to `string`", path, k)
			}
			description = &descriptionString
		}

		var scheme string
		schemeAny, ok := vMap["Scheme"]
		if ok {
			schemeString, ok := schemeAny.(string)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.%s.Scheme\" to `string`", path, k)
			}
			scheme = schemeString
		}

		var format *string = nil
		formatAny, ok := vMap["Format"]
		if ok {
			formatString, ok := formatAny.(string)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.%s.Format\" to `string`", path, k)
			}
			format = &formatString
		}

		var applyToAllRoutes bool
		applyToAllRoutesAny, ok := vMap["ApplyToAllRoutes"]
		if ok {
			applyToAllRoutesString, ok := applyToAllRoutesAny.(bool)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.%s.ApplyToAllRoutes\" to `string`", path, k)
			}
			applyToAllRoutes = applyToAllRoutesString
		}

		ref := self.getRef("", "Auth."+k)
		refHash := hashing.String(ref)

		rootNode, err := getRootNode(path)
		if err != nil {
			return err
		}

		auth := &schemas.Auth{
			Ref:              ref,
			OriginalPath:     self.getPath(fmt.Sprintf("%s.%s", path, k)),
			Name:             k,
			RootNode:         rootNode,
			Description:      description,
			Scheme:           scheme,
			Format:           format,
			ApplyToAllRoutes: applyToAllRoutes,
		}

		stateHash, err := hashing.Struct(auth)
		if err != nil {
			return fmt.Errorf("fail to get state hash for \"%s.%s\"", path, k)
		}

		auth.StateHash = stateHash
		self.schema.Auths.Auths[refHash] = auth
	}

	return nil
}
