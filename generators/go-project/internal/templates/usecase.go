package templates

const UsecaseTempl = `
{{- $pkg := print .DomainSnake "_usecase" -}}
package {{ $pkg }}

import (
{{- range .ImportsUsecase }}
	{{- range . }}
	"{{ . }}"
	{{- end }}
{{ end -}}
)
{{ range $type := .TypesUsecase }}
type {{ $type.GolangType }} struct {
{{- range $prop := $type.MapProps }}
	{{ $prop.Name }}{{ $prop.Spacing1 }} {{ $prop.Type.GetFullTypeName $pkg }}{{ if $prop.Tags }}{{ $prop.Spacing2 }} ` + "`{{ $prop.GetTagsString }}`" + `{{ end }}
{{- end }}
}
{{- end }}

type {{ .Domain }}Usecase interface {
{{- range $method := .MethodsUsecase }}
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
