package templates

const RepositoryTempl = `
{{- $pkg := print .DomainSnake "_repository" -}}
package {{ $pkg }}

import (
{{- range .ImportsRepository }}
	{{- range . }}
	"{{ . }}"
	{{- end }}
{{ end -}}
)
{{ range $type := .TypesRepository }}
type {{ $type.GolangType }} struct {
{{- range $prop := $type.MapProps }}
	{{ $prop.Name }}{{ $prop.Spacing1 }} {{ $prop.Type.GetFullTypeName $pkg }}{{ if $prop.Tags }}{{ $prop.Spacing2 }} ` + "`{{ $prop.GetTagsString }}`" + `{{ end }}
{{- end }}
}
{{- end }}

type {{ .Domain }}Repository interface {
{{- range $method := .MethodsRepository }}
{{- if not $method.OutputTypeName }}
	{{- if not $method.InputTypeName }}
	{{ $method.MethodName }}(ctx context.Context) error
	{{- else }}
	{{ $method.MethodName }}(ctx context.Context, i {{ $method.InputTypeName }}) error
	{{- end }}
{{- else }}
	{{- if not $method.InputTypeName }}
	{{ $method.MethodName }}(ctx context.Context) ({{ $method.OutputTypeName }}, error)
	{{- else }}
	{{ $method.MethodName }}(ctx context.Context, i {{ $method.InputTypeName }}) ({{ $method.OutputTypeName }}, error)
	{{- end }}
{{- end }}
{{- end }}
}
`
