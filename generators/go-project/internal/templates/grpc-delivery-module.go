package templates

const GrpcDeliveryModuleTempl = `
{{- define "method" }}
func (self *{{ .DomainCamel }}Controller) {{ .MethodName }}(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	result, err := self.{{ .DomainCamel }}Usecase.{{ .MethodName }}(ctx)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
{{ end -}}

{{- define "method-with-input" }}
func (self *{{ .DomainCamel }}Controller) {{ .MethodName }}(ctx context.Context, i *pb.{{ .Input.Name }}) (*emptypb.Empty, error) {
	if i == nil {
		return nil, errors.New("input must not be nil")
	}

	err := self.validator.Validate(i)
	if err != nil {
		return nil, err
	}

{{ if .Input.PropsPrepare -}}
	{{- range .Input.PropsPrepare }}
	{{ . }}
	{{- end }}
{{- end }}

	err = self.{{ .DomainCamel }}Usecase.{{ .MethodName }}(ctx, &{{ .DomainSnake }}_usecase.{{ .Input.Name }}{
		{{- range $input := .Input.Props }}
		{{ $input.Name }}:{{ $input.Spacing }} {{ $input.Value }},
		{{- end }}
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
{{ end -}}

{{- define "method-with-output" }}
func (self *{{ .DomainCamel }}Controller) {{ .MethodName }}(ctx context.Context, _ *emptypb.Empty) (*pb.{{ .Output.Name }}, error) {
	result, err := self.{{ .DomainCamel }}Usecase.{{ .MethodName }}(ctx)
	if err != nil {
		return nil, err
	}

{{ if .Output.PropsPrepare -}}
	{{- range .Output.PropsPrepare }}
	{{ . }}
	{{- end }}
{{- end }}

	return &pb.{{ .Output.Name }}{
		{{- range $output := .Output.Props }}
		{{ $output.Name }}:{{ $output.Spacing }} {{ $output.Value }},
		{{- end }}
	}, nil
}
{{ end -}}

{{- define "method-with-input-and-output" }}
func (self *{{ .DomainCamel }}Controller) {{ .MethodName }}(ctx context.Context, i *pb.{{ .Input.Name }}) (*pb.{{ .Output.Name }}, error) {
	if i == nil {
		return nil, errors.New("input must not be nil")
	}

	err := self.validator.Validate(i)
	if err != nil {
		return nil, err
	}

{{ if .Input.PropsPrepare -}}
	{{- range .Input.PropsPrepare }}
	{{ . }}
	{{- end }}
{{- end }}

	result, err := self.{{ .DomainCamel }}Usecase.{{ .MethodName }}(ctx, &{{ .DomainSnake }}_usecase.{{ .Input.Name }}{
		{{- range $input := .Input.Props }}
		{{ $input.Name }}:{{ $input.Spacing }} {{ $input.Value }},
		{{- end }}
	})
	if err != nil {
		return nil, err
	}

{{ if .Output.PropsPrepare -}}
	{{- range .Output.PropsPrepare }}
	{{ . }}
	{{- end }}
{{- end }}

	return &pb.{{ .Output.Name }}{
		{{- range $output := .Output.Props }}
		{{ $output.Name }}:{{ $output.Spacing }} {{ $output.Value }},
		{{- end }}
	}, nil
}
{{ end -}}

package {{ .DomainSnake }}_delivery_grpc

import (
{{- range .ImportsGrpcDelivery }}
	{{- range . }}
	"{{ . }}"
	{{- end }}
{{ end -}}
)

type {{ .DomainCamel }}Controller struct {
	pb.Unimplemented{{ .Domain }}Server

	validator adapters.Validator

	{{ .DomainCamel }}Usecase {{ .DomainSnake }}_usecase.{{ .Domain }}Usecase
}

{{ range $enum := .Enums }}
func convert{{ $enum.GolangName }}ToPb(val models.{{ $enum.GolangName }}) pb.{{ $enum.GolangName }} {
	{{-  range $value := $enum.Values }}
	if val == models.{{ $enum.GolangName }}_{{ $value.Name }} {
		return {{ $value.Idx }}
	}
	{{- end }}

	return 0
}
func convertPbTo{{ $enum.GolangName }}(val pb.{{ $enum.GolangName }}) models.{{ $enum.GolangName }} {
	{{- range $value := $enum.Values }}
	if val == {{ $value.Idx }} {
		return models.{{ $enum.GolangName }}_{{ $value.Name }}
	}
	{{- end }}

	return models.{{ $enum.GolangName }}_{{ with $firstEnum := index $enum.Values 0 }}{{ $firstEnum.Name }}{{ end }}
}
{{- end }}

{{ range .MethodsGrpcDelivery -}}
	{{- if .Input }}
		{{- if .Output }}
		{{- template "method-with-input-and-output" . }}
		{{- else }}
		{{- template "method-with-input" . }}
		{{- end }}
	{{- else }}
		{{- if .Output }}
		{{- template "method-with-output" . }}
		{{- else }}
		{{- template "method" . }}
		{{- end }}
	{{- end }}
{{- end }}

type Add{{ .Domain }}ControllerInput struct {
	Server    *grpc.Server

	Validator adapters.Validator

	{{ .Domain }}Usecase {{ .DomainSnake }}_usecase.{{ .Domain }}Usecase
}

func Add{{ .Domain }}Controller(i *Add{{ .Domain }}ControllerInput) {
	pb.Register{{ .Domain }}Server(i.Server, &{{ .DomainCamel }}Controller{
		validator:{{ .SpacingRelativeToDomainName }}i.Validator,
		{{ .DomainCamel }}Usecase:   i.{{ .Domain }}Usecase,
	})
}
`
