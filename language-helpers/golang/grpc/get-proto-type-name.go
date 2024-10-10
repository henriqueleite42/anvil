package grpc

import (
	"fmt"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *goGrpcParser) GetProtoTypeName(t *schemas.Type) (string, error) {
	if t.Type != schemas.TypeType_Map {
		return "", fmt.Errorf("GetProtoTypeName: type \"%s\"is not a map", t.Name)
	}

	// Adds the domain prefix to avoid duplication
	return self.schema.Domain + t.Name, nil
}
