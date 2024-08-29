package parser

import (
	"fmt"

	"github.com/anvil/anvil/internal/hashing"
	"github.com/anvil/anvil/internal/schema"
)

func (self *anvToAnvpParser) resolveEvent(path string, k string, v any) (string, error) {
	if self.schema.Events == nil {
		self.schema.Events = &schema.Events{}
	}
	if self.schema.Events.Events == nil {
		self.schema.Events.Events = map[string]*schema.Event{}
	}

	originalPath := path + "." + k
	originalPathHash := hashing.String(originalPath)

	_, ok := self.schema.Events.Events[originalPathHash]
	if ok {
		return originalPathHash, nil
	}

	vMap, ok := v.(map[string]any)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s\" to `map[string]any`", path, k)
	}

	// TODO
	_, ok = vMap["$ref"]
	if ok {
		return "", nil
	}

	formatsAny, ok := vMap["Formats"]
	if !ok {
		return "", fmt.Errorf("\"Formats\" is a required property to \"%s.%s\"", path, k)
	}
	formatsArrAny, ok := formatsAny.([]any)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s.Formats\" to `[]any`", path, k)
	}
	formats := []string{}
	for kk, vv := range formatsArrAny {
		formatString, ok := vv.(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Formats.%d\" to `string`", path, k, kk)
		}
		formats = append(formats, formatString)
	}

	eventTypeAny, ok := vMap["Type"]
	if !ok {
		return "", fmt.Errorf("\"Type\" is a required property to \"%s.%s\"", path, k)
	}
	eventTypeHash, err := self.resolveType(AllowedRefs{}, path, k, eventTypeAny)
	if err != nil {
		return "", err
	}

	rootNode, err := getRootNode(path)
	if err != nil {
		return "", err
	}

	event := &schema.Event{
		Name:         k,
		RootNode:     rootNode,
		OriginalPath: originalPath,
		Formats:      formats,
		TypeHash:     eventTypeHash,
	}

	stateHash, err := hashing.Struct(event)
	if err != nil {
		return "", fmt.Errorf("fail to get event \"%s\" state hash", originalPath)
	}

	event.StateHash = stateHash
	self.schema.Events.Events[originalPathHash] = event

	return originalPathHash, nil
}

func (self *anvToAnvpParser) events(file map[string]any) error {
	eventsSchema, ok := file["Events"]
	if !ok {
		return nil
	}

	fullPath := self.getPath("Events")

	eventsMap, ok := eventsSchema.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `map[string]any`", fullPath)
	}

	for k, v := range eventsMap {
		_, err := self.resolveEvent(fullPath, k, v)
		if err != nil {
			return err
		}
	}

	return nil
}
