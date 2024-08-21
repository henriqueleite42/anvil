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
// Usecase and Repository may (and probably will) alter Fields
//
// Doing it at the end ensures that we do it only once
func (self *Parser) stateHashes() error {
	stateHash, err := hashing.Struct(self.schema.Relationships)
	if err != nil {
		return errors.New("fail to get \"Relationships\" state hash")
	}
	self.schema.Relationships.StateHash = stateHash

	stateHash, err = hashing.Struct(self.schema.Imports)
	if err != nil {
		return errors.New("fail to get \"Imports\" state hash")
	}
	self.schema.Imports.StateHash = stateHash

	stateHash, err = hashing.Struct(self.schema.Enums)
	if err != nil {
		return errors.New("fail to get \"Enums\" state hash")
	}
	self.schema.Enums.StateHash = stateHash

	stateHash, err = hashing.Struct(self.schema.Fields)
	if err != nil {
		return errors.New("fail to get \"Fields\" state hash")
	}
	self.schema.Fields.StateHash = stateHash

	stateHash, err = hashing.Struct(self.schema.Types)
	if err != nil {
		return errors.New("fail to get \"Types\" state hash")
	}
	self.schema.Types.StateHash = stateHash

	stateHash, err = hashing.Struct(self.schema.Entities)
	if err != nil {
		return errors.New("fail to get \"Entities\" state hash")
	}
	self.schema.Entities.StateHash = stateHash

	stateHash, err = hashing.Struct(self.schema.Repository)
	if err != nil {
		return errors.New("fail to get \"Repository\" state hash")
	}
	self.schema.Repository.StateHash = stateHash

	stateHash, err = hashing.Struct(self.schema.Usecase)
	if err != nil {
		return errors.New("fail to get \"Usecase\" state hash")
	}
	self.schema.Usecase.StateHash = stateHash

	stateHash, err = hashing.Struct(self.schema.Delivery)
	if err != nil {
		return errors.New("fail to get \"Delivery\" state hash")
	}
	self.schema.Delivery.StateHash = stateHash

	return nil
}

func getRootNode(path string) (string, error) {
	nodes := strings.Split(path, ".")
	if len(nodes) == 0 {
		return "", fmt.Errorf("fail to get root node from \"%s\"", path)
	}
	return nodes[0], nil
}
