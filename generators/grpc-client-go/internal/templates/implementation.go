package templates

const ImplementationTempl = `
{{- $dot := . -}}
package {{ .DomainSnake }}_grpc_client

import (
{{- range .ImportsImplementation }}
	{{- range . }}
	"{{ . }}"
	{{- end }}
{{ end -}}
)

type {{ .DomainCamel }}ApiImplementation struct {
	timeout time.Duration
	conn 		*grpc.ClientConn

	{{ .DomainCamel }}ApiClient pb.{{ .Domain }}ApiClient

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

func (self *{{ .DomainCamel }}ApiImplementation) Close() error {
	return self.conn.Close()
}

{{- range $method := $dot.Methods }}
	{{-  if not $method.Output }}
		{{- if not $method.Input }}
func (self *{{ $dot.DomainCamel }}ApiImplementation) {{ $method.MethodName }}() error {
		{{- else }}
func (self *{{ $dot.DomainCamel }}ApiImplementation) {{ $method.MethodName }}(i {{ $method.Input.GolangType }}) error {
		{{- end }}
	{{- else }}
		{{- if not $method.Input }}
func (self *{{ $dot.DomainCamel }}ApiImplementation) {{ $method.MethodName }}() ({{ $method.Output.GolangType }}, error) {
		{{- else }}
func (self *{{ $dot.DomainCamel }}ApiImplementation) {{ $method.MethodName }}(i {{ $method.Input.GolangType }}) ({{ $method.Output.GolangType }}, error) {
		{{- end }}
	{{- end }}
{{ if $method.Input -}}
{{ range $method.Input.Prepare -}}
{{ . }}
{{ end -}}
{{- end }}
	ctx, cancel := context.WithTimeout(context.Background(), self.timeout)
	defer cancel()
	{{- if not $method.Output }}
		{{- if not $method.Input }}
	_, err := self.{{ $dot.DomainCamel }}ApiClient.{{ $method.MethodName }}(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}
		{{- else }}
	_, err := self.{{ $dot.DomainCamel }}ApiClient.{{ $method.MethodName }}(ctx, {{ $method.Input.Value }})
	if err != nil {
		return err
	}
		{{- end }}
	{{- else }}
		{{- if not $method.Input }}
	result, err := self.{{ $dot.DomainCamel }}ApiClient.{{ $method.MethodName }}(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
		{{- else }}
	result, err := self.{{ $dot.DomainCamel }}ApiClient.{{ $method.MethodName }}(ctx, {{ $method.Input.Value }})
	if err != nil {
		return nil, err
	}
		{{- end }}
	{{- end }}
	{{- if $method.Output }}

{{ range $method.Output.Prepare -}}
{{ . }}
{{ end }}
	return {{ $method.Output.Value }}, nil
	{{- else }}
	return err
	{{- end }}
}
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

	return &{{ .DomainCamel }}ApiImplementation{
		{{ .DomainCamel }}ApiClient: pb.New{{ .Domain }}ApiClient(conn),
		timeout:{{ .SpacingRelativeToDomainName }}timeout,
		conn:   {{ .SpacingRelativeToDomainName }}conn,
	}, nil
}
`
