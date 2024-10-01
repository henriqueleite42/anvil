package templates

const ContractTempl = `package {{ .DomainSnake }}

import (
{{- range .ImportsContract }}
	{{- range . }}
	"{{ . }}"
	{{- end }}
{{ end -}}
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
type {{ $enum.GolangName }} {{ $enum.GolangType }}

const (
{{- range $enumVal := $enum.Values }}
	{{ $enum.GolangName }}_{{ $enumVal.Name }}{{ $enumVal.Spacing }} {{ $enum.GolangName }} = {{ if eq $enum.GolangType "string" }}"{{ $enumVal.Value }}"{{ else }}{{ $enumVal.Value }}{{ end }}
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
type {{ $type.GolangType }} struct {
{{- range $prop := $type.MapProps }}
	{{ $prop.Name }}{{ $prop.Spacing1 }} {{ $prop.GolangType }}
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
	{{ $method.MethodName }}(i *{{ $method.Input.Name }}) error
	{{- end }}
{{- else }}
	{{- if not $method.Input }}
	{{ $method.MethodName }}() (*{{ $method.Output.Name }}, error)
	{{- else }}
	{{ $method.MethodName }}(i *{{ $method.Input.Name }}) (*{{ $method.Output.Name }}, error)
	{{- end }}
{{- end }}
{{- end }}
}
`
