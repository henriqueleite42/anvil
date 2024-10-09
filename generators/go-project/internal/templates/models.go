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

const ModelsTempl = `
{{- $pkg := "models" -}}
package {{ $pkg }}

import (
{{- range .ImportsModels }}
	{{- range . }}
	"{{ . }}"
	{{- end }}
{{ end -}}
)
{{ range $enum := .Enums }}
type {{ $enum.GolangName }} {{ $enum.GolangType }}

const (
{{- range $enumVal := $enum.Values }}
	{{ $enum.GolangName }}_{{ $enumVal.Name }}{{ $enumVal.Spacing }} {{ $enum.GolangName }} = {{ if eq $enum.GolangType "string" }}"{{ $enumVal.Value }}"{{ else }}{{ $enumVal.Value }}{{ end }}
{{- end }}
)
{{ end -}}
{{- range $entity := .Entities }}
type {{ $entity.GolangType }} struct {
{{- range $prop := $entity.MapProps }}
	{{ $prop.Name }}{{ $prop.Spacing1 }} {{ $prop.Type.GetFullTypeName $pkg }}{{ if $prop.Tags }}{{ $prop.Spacing2 }} ` + "`{{ $prop.GetTagsString }}`" + `{{ end }}
{{- end }}
}
{{ end -}}
`
