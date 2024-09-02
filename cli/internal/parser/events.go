package parser

import (
	"fmt"

	"github.com/henriqueleite42/anvil/cli/internal/hashing"
	"github.com/henriqueleite42/anvil/cli/schemas"
)

func (self *anvToAnvpParser) resolveEvent(i *resolveInput) (string, error) {
	if self.schema.Events == nil {
		self.schema.Events = &schemas.Events{}
	}
	if self.schema.Events.Events == nil {
		self.schema.Events.Events = map[string]*schemas.Event{}
	}

	ref := self.getRef(i.ref, "Events."+i.k)
	refHash := hashing.String(ref)

	_, ok := self.schema.Events.Events[refHash]
	if ok {
		return refHash, nil
	}

	vMap, ok := i.v.(map[string]any)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s\" to `map[string]any`", i.path, i.k)
	}

	refAny, ok := vMap["$ref"]
	if ok {
		refString, ok := refAny.(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.$ref\" to `string`", i.path, i.k)
		}
		return hashing.String(refString), nil
	}

	formatsAny, ok := vMap["Formats"]
	if !ok {
		return "", fmt.Errorf("\"Formats\" is a required property to \"%s.%s\"", i.path, i.k)
	}
	formatsArrAny, ok := formatsAny.([]any)
	if !ok {
		return "", fmt.Errorf("fail to parse \"%s.%s.Formats\" to `[]any`", i.path, i.k)
	}
	formats := []string{}
	for kk, vv := range formatsArrAny {
		formatString, ok := vv.(string)
		if !ok {
			return "", fmt.Errorf("fail to parse \"%s.%s.Formats.%d\" to `string`", i.path, i.k, kk)
		}
		formats = append(formats, formatString)
	}

	eventTypeAny, ok := vMap["Type"]
	if !ok {
		return "", fmt.Errorf("\"Type\" is a required property to \"%s.%s\"", i.path, i.k)
	}
	eventTypeHash, err := self.resolveType(&resolveInput{
		path: i.path,
		ref:  "Events",
		k:    i.k,
		v:    eventTypeAny,
	})
	if err != nil {
		return "", err
	}

	rootNode, err := getRootNode(i.path)
	if err != nil {
		return "", err
	}

	event := &schemas.Event{
		Ref:          ref,
		OriginalPath: fmt.Sprintf("%s.%s", i.path, i.k),
		Name:         i.k,
		RootNode:     rootNode,
		Formats:      formats,
		TypeHash:     eventTypeHash,
	}

	stateHash, err := hashing.Struct(event)
	if err != nil {
		return "", fmt.Errorf("fail to get event \"%s.%s\" state hash", i.path, i.k)
	}

	event.StateHash = stateHash
	self.schema.Events.Events[refHash] = event

	return refHash, nil
}

func (self *anvToAnvpParser) events(file map[string]any) error {
	eventsSchema, ok := file["Events"]
	if !ok {
		return nil
	}

	path := "Events"

	eventsMap, ok := eventsSchema.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `map[string]any`", path)
	}

	for k, v := range eventsMap {
		_, err := self.resolveEvent(&resolveInput{
			path: path,
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
