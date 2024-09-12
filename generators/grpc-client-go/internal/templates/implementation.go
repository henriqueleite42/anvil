package templates

const ImplementationTempl = `package {{ .DomainSnake }}

import (
{{- range .ImportsImplementation }}
{{ . }}
{{- end }}
)

type {{ .Domain }}ApiImplementation struct {
	timeout time.Duration

	{{ .DomainCamel }}Client pb.{{ .Domain }}Client

	conn *grpc.ClientConn
}
{{ range $enum := .Enums }}
func convert{{ $enum.Name }}ToPb(val {{ $enum.Name }}) pb.{{ $enum.Name }} {
	{{-  range $value := $enum.Values }}
	if val == {{ $enum.Name }}_{{ $value.Name }} {
		return {{ $value.Idx }}
	}
	{{- end }}

	return 0
}
func convertPbTo{{ $enum.Name }}(val pb.{{ $enum.Name }}) {{ $enum.Name }} {
	{{- range $value := $enum.Values }}
	if val == {{ $value.Idx }} {
		return {{ $enum.Name }}_{{ $value.Name }}
	}
	{{- end }}

	return {{ $enum.Name }}_{{ with $firstEnum := index $enum.Values 0 }}{{ $firstEnum.Name }}{{ end }}
}
{{- end }}

func (self *{{ .Domain }}ApiImplementation) Close() error {
	return self.conn.Close()
}

{{- with $originalInput := . }}
{{- range $method := $originalInput.Methods }}
	{{-  if not $method.Output }}
		{{- if not $method.Input }}
func (self *{{ $originalInput.Domain }}ApiImplementation) {{ $method.MethodName }}() error {
		{{- else }}
func (self *{{ $originalInput.Domain }}ApiImplementation) {{ $method.MethodName }}(i *{{ $method.MethodName }}Input) error {
		{{- end }}
	{{- else }}
		{{- if not $method.Input }}
func (self *{{ $originalInput.Domain }}ApiImplementation) {{ $method.MethodName }}() (*{{ $method.MethodName }}Output, error) {
		{{- else }}
func (self *{{ $originalInput.Domain }}ApiImplementation) {{ $method.MethodName }}(i *{{ $method.MethodName }}Input) (*{{ $method.MethodName }}Output, error) {
		{{- end }}
	{{- end }}
	{{- if $method.Input }}
	if i == nil {
		return errors.New("input must not be nil")
	}
{{ range $method.InputPropsPrepare }}
{{ . }}
{{- end }}
{{ end }}{{/* NECESSARY IT TO BE THIS WAY!!!! */}}
	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()
	{{- if not $method.Output }}
		{{- if not $method.Input }}
	_, err := self.{{ $originalInput.DomainCamel }}Client.{{ $method.MethodName }}(ctx, nil)
		{{- else }}
	_, err := self.{{ $originalInput.DomainCamel }}Client.{{ $method.MethodName }}(ctx, &pb.{{ $method.MethodName }}Input{
		{{- range $input := $method.Input }}
		{{ $input.Name }}:{{ $input.Spacing }} {{ $input.Value }},
		{{- end }}
	})
		{{- end }}
	{{- else }}
		{{- if not $method.Input }}
	result, err := self.{{ $originalInput.DomainCamel }}Client.{{ $method.MethodName }}(ctx, nil)
		{{- else }}
	result, err := self.{{ $originalInput.DomainCamel }}Client.{{ $method.MethodName }}(ctx, &pb.{{ $method.MethodName }}Input{
		{{- range $input := $method.Input }}
		{{ $input.Name }}:{{ $input.Spacing }} {{ $input.Value }},
		{{- end }}
	})
		{{- end }}
	{{- end }}
	{{- if $method.Output }}

{{ $method.OutputPropsPrepare }}

	return &{{ $method.MethodName }}Output{
		{{- range $output := $method.Output }}
	{{ $output.Name }}:{{ $output.Spacing }} {{ $output.Value }},
		{{- end }}
	},err
	{{- else }}
	return err
	{{- end }}
}
{{- end }}
{{- end }}

func New{{ .Domain }}Api(i *{{ .Domain }}ApiInput) ({{ .Domain }}Api, error) {
	if i == nil {
		return nil, errors.New("\"New{{ .Domain }}Api\": input must not be nil")
	}

	if i.Addr == "" {
		return nil, errors.New("\"Addr\" is required")
	}

	conn, err := grpc.NewClient(i.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	var timeout time.Duration
	if i.Timeout == 0 {
		timeout = time.Second * 5
	} else {
		timeout = i.Timeout
	}

	return &{{ .Domain }}ApiImplementation{
		{{ .DomainCamel }}Client: pb.New{{ .Domain }}Client(conn),
		timeout:{{ .SpacingRelativeToDomainName }}timeout,
		conn:   {{ .SpacingRelativeToDomainName }}conn,
	}, nil
}
`
