package templates

const ContractTempl = `
{{- $dot := . -}}
package {{ .DomainSnake }}_grpc_client

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
	{{ $prop.Name }}{{ $prop.Spacing1 }} {{ $prop.Type.GetFullTypeName $dot.DomainSnake }}
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
	{{ $method.MethodName }}(i {{ $method.Input.GolangType }}) error
	{{- end }}
{{- else }}
	{{- if not $method.Input }}
	{{ $method.MethodName }}() ({{ $method.Output.GolangType }}, error)
	{{- else }}
	{{ $method.MethodName }}(i {{ $method.Input.GolangType }}) ({{ $method.Output.GolangType }}, error)
	{{- end }}
{{- end }}
{{- end }}
}
`
