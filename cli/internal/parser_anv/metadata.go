package parser_anv

import (
	"fmt"

	"github.com/henriqueleite42/anvil/cli/schemas"
)

func (self *anvToAnvpParser) metadata(file map[string]any) error {
	metadataSchema, ok := file["Metadata"]
	if !ok {
		return nil
	}

	fullPath := self.getPath("Metadata")

	valMap, ok := metadataSchema.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `map[string]any`", fullPath)
	}

	metadata := &schemas.Metadata{}

	description, ok := valMap["Description"]
	if ok {
		valString, ok := description.(string)
		if !ok {
			return fmt.Errorf("fail to parse \"%s.Description\" to `string`", fullPath)
		}
		metadata.Description = valString
	}

	serversAny, ok := valMap["Servers"]
	if ok {
		valMap, ok := serversAny.(map[string]any)
		if !ok {
			return fmt.Errorf("fail to parse \"%s.Servers\" to `map[string]any`", fullPath)
		}

		servers := map[string]*schemas.MetadataServers{}

		for k, v := range valMap {
			vMap, ok := v.(map[string]any)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.Servers.%s\" to `map[string]any`", fullPath, k)
			}

			urlAny, ok := vMap["Url"]
			if !ok {
				return fmt.Errorf("\"Url\" is a required property to \"%s.Servers.%s\"", fullPath, k)
			}
			urlString, ok := urlAny.(string)
			if !ok {
				return fmt.Errorf("fail to parse \"%s.Servers.%s.Url\" to `string`", fullPath, k)
			}

			servers[k] = &schemas.MetadataServers{
				Url: urlString,
			}
		}

		metadata.Servers = servers
	}

	self.schema.Metadata = metadata

	return nil
}
