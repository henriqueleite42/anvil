package templates

type TemplEnumValue struct {
	Idx     int
	Name    string
	Spacing string
	Value   string
}

type TemplEnum struct {
	Name   string
	Type   string
	Values []*TemplEnumValue
}

const ModelsTempl = `package {{ .DomainSnake }}_models

import (
{{- range .ImportsModels }}
{{ . }}
{{- end }}
)
{{ range $enum := .Enums }}
type {{ $enum.Name }} {{ $enum.Type }}

const (
{{- range $enumVal := $enum.Values }}
	{{ $enum.Name }}_{{ $enumVal.Name }}{{ $enumVal.Spacing }} {{ $enum.Name }} = {{ if eq $enum.Type "string" }}"{{ $enumVal.Value }}"{{ else }}{{ $enumVal.Value }}{{ end }}
{{- end }}
)
{{ end -}}
{{- range $entity := .Entities }}
type {{ $entity.Name }} struct {
{{- range $prop := $entity.Props }}
	{{ $prop.Name }}{{ $prop.Spacing1 }} {{ $prop.Type }}{{ if $prop.Tags }}{{ $prop.Spacing2 }} {{ $prop.Tags }}{{ end }}
{{- end }}
}
{{ end -}}
`
