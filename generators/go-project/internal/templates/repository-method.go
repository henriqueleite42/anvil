package templates

type RepositoryMethodTemplInput struct {
	Domain         string
	DomainSnake    string
	MethodName     string
	InputTypeName  string
	OutputTypeName string
}

const RepositoryMethodTempl = `package {{ .DomainSnake }}_repository

import (
	"context"
	"errors"
)

{{ if not .OutputTypeName -}}
	{{- if not .InputTypeName -}}

func (self *{{ .Domain }}RepositoryImplementation) {{ .MethodName }}(ctx context.Context) error {
	return errors.New("\"{{ .MethodName }}\": unimplemented")
}

	{{- else -}}

func (self *{{ .Domain }}RepositoryImplementation) {{ .MethodName }}(ctx context.Context, i {{ .InputTypeName }}) error {
	return errors.New("\"{{ .MethodName }}\": unimplemented")
}

	{{- end -}}
{{- else -}}
	{{- if not .InputTypeName -}}

func (self *{{ .Domain }}RepositoryImplementation) {{ .MethodName }}(ctx context.Context) ({{ .OutputTypeName }}, error) {
	return nil, errors.New("\"{{ .MethodName }}\": unimplemented")
}

	{{- else -}}

func (self *{{ .Domain }}RepositoryImplementation) {{ .MethodName }}(ctx context.Context, i {{ .InputTypeName }}) ({{ .OutputTypeName }}, error) {
	return nil, errors.New("\"{{ .MethodName }}\": unimplemented")
}

	{{- end }}
{{- end }}
`
