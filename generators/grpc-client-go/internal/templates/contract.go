package templates

const ContractTempl = `package {{ .DomainSnake }}

import (
{{- range .ImportsContract }}
{{ . }}
{{- end }}
)

type {{ .Domain }}ApiInput struct {
	Addr    string
	Timeout time.Duration
}

{{- /*
------------------
-
- ENUMS
-
------------------
*/}}
{{ range $enum := .Enums }}
type {{ $enum.Name }} {{ $enum.Type }}
const (
{{- range $enumVal := $enum.Values }}
	{{ $enum.Name }}_{{ $enumVal.Name }}{{ $enumVal.Spacing }} {{ $enum.Name }} = {{ if eq $enum.Type "string" }}"{{ $enumVal.Value }}"{{ else }}{{ $enumVal.Value }}{{ end }}
{{- end }}
)
{{- end }}
{{- /*
	------------------
	-
	- TYPES
	-
	------------------
*/}}
{{ range $type := .Types }}
type {{ $type.Name }} struct {
{{- range $prop := $type.Props }}
	{{ $prop.Name }}{{ $prop.Spacing }} {{ $prop.Type }}
{{- end }}
}
{{- end }}

type {{ .Domain }}Api interface {
	Close() error
{{ range $method := .Methods }}
{{- if not $method.Output }}
	{{- if not $method.Input }}
	{{ $method.MethodName }}() error
	{{- else }}
	{{ $method.MethodName }}(i *{{ $method.InputTypeName }}) error
	{{- end }}
{{- else }}
	{{- if not $method.Input }}
	{{ $method.MethodName }}() (*{{ $method.OutputTypeName }}, error)
	{{- else }}
	{{ $method.MethodName }}(i *{{ $method.InputTypeName }}) (*{{ $method.OutputTypeName }}, error)
	{{- end }}
{{- end }}
{{- end }}
}
`
