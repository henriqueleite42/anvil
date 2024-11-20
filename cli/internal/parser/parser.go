package parser

import (
	"github.com/henriqueleite42/anvil/cli/internal/files"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

type anvToAnvpParser struct {
	schema *schemas.AnvpSchema
}

type resolveInput struct {
	namePrefix string // Internal use. Correctly parse child map types with the prefix of their parent.

	curDomain string
	path      string // Original path
	ref       string // Ref until now
	k         string // Key being resolved, usually the type name, but if it's an child type, it's only part of the type's name
	v         any    // Value, type specification
}

func (self *anvToAnvpParser) parseCommon(filesAny map[string]map[string]any) error {
	// ---------------------
	//
	// Needs to parse one by one!
	// First, all metadata from all domains,
	// them, all the auth from all domains,
	// and so on
	// It's done this way, because domains can be dependents
	// of other domains. Ex:
	// A Type from domain A is dependent of an Enum of domain B
	//
	// ---------------------
	for curDomain, file := range filesAny {
		err := self.resolveEntitiesMetadata(curDomain, file)
		if err != nil {
			return err
		}
	}

	for curDomain, file := range filesAny {
		err := self.auth(curDomain, file)
		if err != nil {
			return err
		}
	}

	for curDomain, file := range filesAny {
		err := self.enums(curDomain, file)
		if err != nil {
			return err
		}
	}

	for curDomain, file := range filesAny {
		err := self.types(curDomain, file)
		if err != nil {
			return err
		}
	}

	for curDomain, file := range filesAny {
		err := self.entities(curDomain, file)
		if err != nil {
			return err
		}
	}

	for curDomain, file := range filesAny {
		err := self.events(curDomain, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *anvToAnvpParser) parse(curDomain string, file map[string]any) error {
	err := self.repository(curDomain, file)
	if err != nil {
		return err
	}

	err = self.usecase(curDomain, file)
	if err != nil {
		return err
	}

	err = self.delivery(curDomain, file)
	if err != nil {
		return err
	}

	return nil
}

func ParseAnvToAnvp(uris []string) (*schemas.AnvpSchema, error) {
	parser := &anvToAnvpParser{
		schema: &schemas.AnvpSchema{},
	}

	filesAny := make(map[string]map[string]any, len(uris))

	for _, uri := range uris {
		file, err := files.ReadAnvFile(uri)
		if err != nil {
			return nil, err
		}

		curDomain, err := parser.domain(uri, file)
		if err != nil {
			return nil, err
		}

		filesAny[curDomain] = file
	}

	err := parser.parseCommon(filesAny)
	if err != nil {
		return nil, err
	}

	for curDomain, file := range filesAny {
		err := parser.parse(curDomain, file)
		if err != nil {
			return nil, err
		}

		err = parser.stateHashes()
		if err != nil {
			return nil, err
		}
	}

	return parser.schema, nil
}
