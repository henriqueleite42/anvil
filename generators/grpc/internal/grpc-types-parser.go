package internal

import (
	"github.com/henriqueleite42/anvil/generators/grpc/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/imports"
)

// One types parser per domain
type grpcTypesParser struct {
	imports  imports.ImportsManager
	methods  []*templates.ProtofileTemplInputMethod
	enums    map[string]*templates.ProtofileTemplInputEnum
	types    []*templates.ProtofileTemplInputType
	events   []*templates.ProtofileTemplInputType
	entities []*templates.ProtofileTemplInputType
}
