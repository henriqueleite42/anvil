package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/henriqueleite42/anvil/language-helpers/golang/formatter"
	"github.com/henriqueleite42/anvil/language-helpers/golang/hashing"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

type SortedByOrder struct {
	Order int
	Key   string
}

//	Add the state hash to all the properties
//
// We do it separately instead of doing it in the parsers because
// They may get modified by other parsers. Example:
// Usecase and Repository may (and probably will) alter Types
//
// Doing it at the end ensures that we do it only once
func (self *anvToAnvpParser) stateHashes() error {
	if self.schema.Auths != nil {
		stateHash, err := hashing.Struct(self.schema.Auths)
		if err != nil {
			return errors.New("fail to get \"Auths\" state hash")
		}
		self.schema.Auths.StateHash = stateHash
	}

	if self.schema.Enums != nil {
		stateHash, err := hashing.Struct(self.schema.Enums)
		if err != nil {
			return errors.New("fail to get \"Enums\" state hash")
		}
		self.schema.Enums.StateHash = stateHash
	}

	if self.schema.Types != nil {
		stateHash, err := hashing.Struct(self.schema.Types)
		if err != nil {
			return errors.New("fail to get \"Types\" state hash")
		}
		self.schema.Types.StateHash = stateHash
	}

	if self.schema.Events != nil {
		stateHash, err := hashing.Struct(self.schema.Events)
		if err != nil {
			return errors.New("fail to get \"Events\" state hash")
		}
		self.schema.Events.StateHash = stateHash
	}

	if self.schema.Entities != nil {
		stateHash, err := hashing.Struct(self.schema.Entities)
		if err != nil {
			return errors.New("fail to get \"Entities\" state hash")
		}
		self.schema.Entities.StateHash = stateHash
	}

	if self.schema.Repositories != nil {
		stateHash, err := hashing.Struct(self.schema.Repositories)
		if err != nil {
			return errors.New("fail to get \"Repositories\" state hash")
		}
		self.schema.Repositories.StateHash = stateHash
	}

	if self.schema.Usecases != nil {
		stateHash, err := hashing.Struct(self.schema.Repositories)
		if err != nil {
			return errors.New("fail to get \"Repositories\" state hash")
		}
		self.schema.Repositories.StateHash = stateHash
	}

	if self.schema.Deliveries != nil {
		stateHash, err := hashing.Struct(self.schema.Deliveries)
		if err != nil {
			return errors.New("fail to get \"Deliveries\" state hash")
		}
		self.schema.Deliveries.StateHash = stateHash
	}

	return nil
}

type GetRefInput struct {
	SchemaProperty string
	Name           string
	FieldName      string // Optional, not using pointer to facilitate use
	NestedRef      string // Optional, not using pointer to facilitate use
}

func (self *anvToAnvpParser) getRef(curDomain string, ref string) string {
	return self.getDeepRef(curDomain, "", ref)
}

func (self *anvToAnvpParser) getDeepRef(curDomain string, parentRef string, ref string) string {
	domainPref := curDomain + "."

	refWithoutDomain := strings.TrimPrefix(ref, domainPref)

	if parentRef == "" {
		return domainPref + refWithoutDomain
	}

	parentRefWithoutDomain := strings.TrimPrefix(parentRef, domainPref)

	return domainPref + parentRefWithoutDomain + "." + ref
}

func (self *anvToAnvpParser) anvRefToAnvpRef(
	curDomain string,
	ref string,
) string {
	if strings.HasPrefix(ref, "Entities") ||
		strings.HasPrefix(ref, "Types") ||
		strings.HasPrefix(ref, "Enums") ||
		strings.HasPrefix(ref, "Auths") ||
		strings.HasPrefix(ref, "Events") {
		return curDomain + "." + ref
	}

	return ref
}

func getRootNode(path string) (string, error) {
	nodes := strings.Split(path, ".")
	if len(nodes) < 2 {
		return "", fmt.Errorf("fail to get root node from \"%s\"", path)
	}
	return nodes[1], nil
}

func (self *anvToAnvpParser) formatToEntitiesNamingCase(str string) string {
	if self.schema.Entities == nil {
		return str
	}
	if self.schema.Entities.Metadata == nil {
		return str
	}
	if self.schema.Entities.Metadata.NamingCase == nil {
		return str
	}

	switch *self.schema.Entities.Metadata.NamingCase {
	case schemas.NamingCase_Camel:
		return formatter.PascalToCamel(str)
	case schemas.NamingCase_Snake:
		return formatter.PascalToSnake(str)
	}

	return str
}
