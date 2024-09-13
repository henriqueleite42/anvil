package templates

const RepositoryTempl = `package {{ .DomainSnake }}_repository

import (
{{- range .ImportsRepository }}
{{ . }}
{{- end }}
)
{{ range $type := .TypesRepository }}
type {{ $type.Name }} struct {
{{- range $prop := $type.Props }}
	{{ $prop.Name }}{{ $prop.Spacing1 }} {{ $prop.Type }}{{ $prop.Spacing2 }} {{ $prop.Tags }}
{{- end }}
}
{{- end }}

type {{ .Domain }}Repository interface {
{{- range $method := .MethodsRepository }}
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
