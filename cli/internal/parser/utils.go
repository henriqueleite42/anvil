package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/anvil/anvil/internal/hashing"
)

//	Add the state hash to all the properties
//
// We do it separately instead of doing it in the parsers because
// They may get modified by other parsers. Example:
// Usecase and Repository may (and probably will) alter Types
//
// Doing it at the end ensures that we do it only once
func (self *anvToAnvpParser) stateHashes() error {
	if self.schema.Relationships != nil {
		stateHash, err := hashing.Struct(self.schema.Relationships)
		if err != nil {
			return errors.New("fail to get \"Relationships\" state hash")
		}
		self.schema.Relationships.StateHash = stateHash
	}

	if self.schema.Imports != nil {
		stateHash, err := hashing.Struct(self.schema.Imports)
		if err != nil {
			return errors.New("fail to get \"Imports\" state hash")
		}
		self.schema.Imports.StateHash = stateHash
	}

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

	if self.schema.Repository != nil {
		stateHash, err := hashing.Struct(self.schema.Repository)
		if err != nil {
			return errors.New("fail to get \"Repository\" state hash")
		}
		self.schema.Repository.StateHash = stateHash
	}

	if self.schema.Usecase != nil {
		stateHash, err := hashing.Struct(self.schema.Usecase)
		if err != nil {
			return errors.New("fail to get \"Usecase\" state hash")
		}
		self.schema.Usecase.StateHash = stateHash
	}

	if self.schema.Delivery != nil {
		stateHash, err := hashing.Struct(self.schema.Delivery)
		if err != nil {
			return errors.New("fail to get \"Delivery\" state hash")
		}
		self.schema.Delivery.StateHash = stateHash
	}

	return nil
}

func (self *anvToAnvpParser) getPath(path string) string {
	if self.baseRef == "" {
		return path
	}

	return self.baseRef + "." + path
}

func getRootNode(path string) (string, error) {
	nodes := strings.Split(path, ".")
	if len(nodes) == 0 {
		return "", fmt.Errorf("fail to get root node from \"%s\"", path)
	}
	return nodes[0], nil
}

type GetRefInput struct {
	SchemaProperty string
	Name           string
	FieldName      string // Optional, not using pointer to facilitate use
	NestedRef      string // Optional, not using pointer to facilitate use
}

func (self *anvToAnvpParser) getRef(parentRef string, ref string) string {
	if parentRef == "" {
		return ref
	}

	return parentRef + "." + ref

}

func (self *anvToAnvpParser) getRefHash(ref string) string {
	fullRef := self.baseRef + "." + ref
	return hashing.String(fullRef)
}
