package parser

import (
	"errors"
	"fmt"

	"github.com/anvil/anvil/internal/schema"
)

func (self *Parser) metadata(file map[string]any) error {
	metadataSchema, ok := file["Metadata"]
	if !ok {
		return nil
	}

	valMap, ok := metadataSchema.(map[string]any)
	if !ok {
		return errors.New("fail to parse \"Metadata\" to `map[string]any`")
	}

	metadata := &schema.Metadata{}

	description, ok := valMap["Description"]
	if ok {
		valString, ok := description.(string)
		if !ok {
			return errors.New("fail to parse \"Metadata.Description\" to `string`")
		}
		metadata.Description = valString
	}

	serversAny, ok := valMap["Servers"]
	if ok {
		valMap, ok := serversAny.(map[string]any)
		if !ok {
			return errors.New("fail to parse \"Metadata.Servers\" to `map[string]any`")
		}

		servers := map[string]*schema.MetadataServers{}

		for k, v := range valMap {
			vMap, ok := v.(map[string]any)
			if !ok {
				return fmt.Errorf("fail to parse \"Metadata.Servers.%s\" to `map[string]any`", k)
			}

			urlAny, ok := vMap["Url"]
			if !ok {
				return fmt.Errorf("\"Url\" is a required property to \"Metadata.Servers.%s\"", k)
			}
			urlString, ok := urlAny.(string)
			if !ok {
				return fmt.Errorf("fail to parse \"Metadata.Servers.%s.Url\" to `string`", k)
			}

			servers[k] = &schema.MetadataServers{
				Url: urlString,
			}
		}

		metadata.Servers = servers
	}

	self.schema.Metadata = metadata

	return nil
}
