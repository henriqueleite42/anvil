package parser

import (
	"github.com/henriqueleite42/anvil/cli/schemas"
	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
)

type Parser struct {
	ModelsPath    string
	ModelsPkgName string
	Schema        *schemas.Schema

	ImportsModels     map[string]bool
	ImportsRepository map[string]bool
	ImportsUsecase    map[string]bool
	Enums             map[string]*templates.TemplEnum
	Entities          []*templates.TemplType
	TypesRepository   []*templates.TemplType
	TypesUsecase      []*templates.TemplType
	MethodsRepository []*templates.TemplMethod
	MethodsUsecase    []*templates.TemplMethod
}

type Kind string

const (
	Kind_Entity     Kind = "ENTITY"
	Kind_Event      Kind = "EVENT"
	Kind_Repository Kind = "REPOSITORY"
	Kind_Usecase    Kind = "USECASE"
)
