package templates

import "github.com/henriqueleite42/anvil/language-helpers/golang/schemas"

type HclTemplInputEntityColumn struct {
	Name          string
	Type          string
	Default       *string
	Optional      bool
	AutoIncrement bool
	Order         int
}

type HclTemplInputEntityPrimaryColumn struct {
	DbName string
	Order  int
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

const HclTempl = `schema "public" {}

{{ range .Enums -}}
enum "{{ .DbName }}" {
	schema = schema.public
	values = [
		{{- range .Values }}
		"{{ .Value }}",
		{{- end }}
	]
}
{{ end }}
{{ range .Entities }}
table "{{ .DbName }}" {
	schema = schema.public
	{{- range .Columns }}
	column "{{ .Name }}" {
		type = {{ .Type }}
		{{- if .AutoIncrement }}
    identity {
			generated = ALWAYS
			start = 0
			increment = 1
    }
		{{- end }}
		{{- if .Optional }}
		null = true
		{{- end }}
		{{- if .Default }}
		default = {{ .Default }}
		{{- end }}
	}
	{{- end}}
	primary_key {
		columns = [
			{{- range .PrimaryKeys }}
			{{ .DbName }},
			{{- end}}
		]
	}
	{{- range .Indexes }}
	index "{{ .Name }}" {
		columns = [
			{{- range .Columns }}
			{{ . }}
			{{- end }}
		]
		{{- if .Unique }}
		unique = true
		{{- end }}
	}
	{{- end }}
	{{- range .ForeignKeys }}
	{{- $refTable := .RefTable }}
	foreign_key "{{ .Name }}" {
		columns = [
			{{- range .Columns }}
			{{ . }}
			{{- end }}
		]
		ref_columns = [
			{{- range .RefColumns }}
			{{ . }}
			{{- end }}
		]
		{{- if .OnUpdate }}
		on_update = {{ .OnUpdate }}
		{{- end }}
		{{- if .OnDelete }}
		on_delete = {{ .OnDelete }}
		{{- end }}
	}
	{{- end }}
}
{{- end }}
`
