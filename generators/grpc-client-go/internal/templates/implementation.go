package templates

const ImplementationTempl = `package {{ .DomainSnake }}

import (
{{- range .ImportsImplementation }}
	{{- range . }}
	"{{ . }}"
	{{- end }}
{{ end -}}
)

type {{ .Domain }}ApiImplementation struct {
	timeout time.Duration

	{{ .DomainCamel }}Client pb.{{ .Domain }}Client

	conn *grpc.ClientConn
}
{{ range $enum := .Enums }}
func convert{{ $enum.GolangName }}ToPb(val {{ $enum.GolangName }}) pb.{{ $enum.GolangName }} {
	{{-  range $value := $enum.Values }}
	if val == {{ $enum.GolangName }}_{{ $value.Name }} {
		return {{ $value.Idx }}
	}
	{{- end }}

	return 0
}
func convertPbTo{{ $enum.GolangName }}(val pb.{{ $enum.GolangName }}) {{ $enum.GolangName }} {
	{{- range $value := $enum.Values }}
	if val == {{ $value.Idx }} {
		return {{ $enum.GolangName }}_{{ $value.Name }}
	}
	{{- end }}

	return {{ $enum.GolangName }}_{{ with $firstEnum := index $enum.Values 0 }}{{ $firstEnum.Name }}{{ end }}
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
func (self *{{ $originalInput.Domain }}ApiImplementation) {{ $method.MethodName }}(i *{{ $method.Input.Name }}) error {
	if i == nil {
		return errors.New("input must not be nil")
	}
		{{- end }}
	{{- else }}
		{{- if not $method.Input }}
func (self *{{ $originalInput.Domain }}ApiImplementation) {{ $method.MethodName }}() (*{{ $method.Output.Name }}, error) {
		{{- else }}
func (self *{{ $originalInput.Domain }}ApiImplementation) {{ $method.MethodName }}(i *{{ $method.Input.Name }}) (*{{ $method.Output.Name }}, error) {
	if i == nil {
		return nil, errors.New("input must not be nil")
	}
		{{- end }}
	{{- end }}

{{ if $method.Input -}}
{{ range $method.Input.PropsPrepare -}}
{{ . }}
{{ end -}}
{{- end }}
	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()
	{{- if not $method.Output }}
		{{- if not $method.Input }}
	_, err := self.{{ $originalInput.DomainCamel }}Client.{{ $method.MethodName }}(ctx, &emptypb.Empty{})
		{{- else }}
	_, err := self.{{ $originalInput.DomainCamel }}Client.{{ $method.MethodName }}(ctx, &pb.{{ $method.Input.Name }}{
		{{- range $input := $method.Input.Props }}
		{{ $input.Name }}:{{ $input.Spacing }} {{ $input.Value }},
		{{- end }}
	})
		{{- end }}
	{{- else }}
		{{- if not $method.Input }}
	result, err := self.{{ $originalInput.DomainCamel }}Client.{{ $method.MethodName }}(ctx, &emptypb.Empty{})
		{{- else }}
	result, err := self.{{ $originalInput.DomainCamel }}Client.{{ $method.MethodName }}(ctx, &pb.{{ $method.Input.Name }}{
		{{- range $input := $method.Input.Props }}
		{{ $input.Name }}:{{ $input.Spacing }} {{ $input.Value }},
		{{- end }}
	})
		{{- end }}
	{{- end }}
	{{- if $method.Output }}

{{ range $method.Output.PropsPrepare -}}
{{ . }}
{{ end }}
	return &{{ $method.Output.Name }}{
		{{- range $output := $method.Output.Props }}
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
