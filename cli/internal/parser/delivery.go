package parser

import (
	"fmt"

	"github.com/henriqueleite42/anvil/cli/schemas"
)

func (self *anvToAnvpParser) delivery(file map[string]any) error {
	path := self.getPath("Delivery")

	deliveryAny, ok := file["Delivery"]
	if ok {
		return nil
	}

	_, ok = deliveryAny.(map[string]any)
	if !ok {
		return fmt.Errorf("fail to parse \"%s\" to `map[string]any`", path)
	}

	// TODO parse delivery

	// TODO parse grpc

	// TODO parse http

	// TODO parse queue

	self.schema.Delivery = &schemas.Delivery{}

	return nil
}
