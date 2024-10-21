package templates

import (
	_ "embed"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

type HclTemplInputEntityColumn struct {
	Name          string
	Type          string
	Default       *string
	Optional      bool
	AutoIncrement bool
	Order         uint
}

type HclTemplInputEntityPrimaryColumn struct {
	DbName string
	Order  uint
}

type HclTemplInputEntityForeignKey struct {
	RefTable   string
	Name       string
	Columns    []string
	RefColumns []string
	OnUpdate   *string
	OnDelete   *string
}

type HclTemplInputEntityIndex struct {
	Name    string
	Columns []string
	Unique  bool
}

type HclTemplInputEntity struct {
	DbName      string
	Columns     []*HclTemplInputEntityColumn
	PrimaryKeys []*HclTemplInputEntityPrimaryColumn
	Indexes     []*HclTemplInputEntityIndex
	ForeignKeys []*HclTemplInputEntityForeignKey
}

type HclTemplInput struct {
	Enums    []*schemas.Enum
	Entities []*HclTemplInputEntity
}

//go:embed hcl.txt
var HclTempl string
