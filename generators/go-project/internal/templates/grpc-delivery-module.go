package templates

const GrpcDeliveryModuleTempl = `
{{- define "logger" }}	logger := self.logger.With().
		Str("dmn", "{{ .Domain }}").
		Str("mtd", "{{ .MethodName }}").
		Str("reqId", xid.New().String()).
		Logger()
	logger.Trace().Msg("start")
{{- end -}}

{{- define "input" }}	logger.Trace().Msg("check i == nil")
	if i == nil {
		logger.Error().Msg("input is nil")
		return nil, errors.New("input must not be nil")
	}

	logger.Trace().Msg("validate i")
	err := self.validator.Validate(i)
	if err != nil {
		logger.Info().Err(err).Msg("invalid i")
		return nil, err
	}

{{ if .Input.PropsPrepare -}}
	logger.Trace().Msg("start input props prepare")
	{{- range .Input.PropsPrepare }}
{{ . }}
	{{- end }}
	logger.Trace().Msg("end input props prepare")
{{- end }}

	logger.Trace().Msg("build uscI")
	uscI := &{{ .DomainSnake }}_usecase.{{ .Input.Name }}{
		{{- range $input := .Input.Props }}
		{{ $input.Name }}:{{ $input.Spacing }} {{ $input.Value }},
		{{- end }}
	}
	logger.Debug().Interface("uscI", uscI).Msg("usecase input")
{{- end -}}

{{- define "output" }}	logger.Debug().Interface("result", result).Msg("usecase output")
{{ if .Output.PropsPrepare }}
	logger.Trace().Msg("start output props prepare")
	{{- range .Output.PropsPrepare }}
{{ . }}
	{{- end }}
	logger.Trace().Msg("end output props prepare")
{{- end }}

	logger.Trace().Msg("end")
	return &pb.{{ .Output.Name }}{
		{{- range $output := .Output.Props }}
		{{ $output.Name }}:{{ $output.Spacing }} {{ $output.Value }},
		{{- end }}
	}, nil
{{- end -}}

{{- define "method" }}
func (self *{{ .DomainCamel }}Controller) {{ .MethodName }}(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
{{ template "logger" . }}

	logger.Trace().Msg("create reqCtx")
	reqCtx := context.WithValue(ctx, "logger", logger)

	logger.Trace().Msg("call usecase")
	err := self.{{ .DomainCamel }}Usecase.{{ .MethodName }}(reqCtx)
	if err != nil {
		logger.Warn().Err(err).Msg("usecase err")
		return nil, err
	}

	logger.Trace().Msg("end")
	return &emptypb.Empty{}, nil
}
{{- end -}}

{{- define "method-with-input" }}
func (self *{{ .DomainCamel }}Controller) {{ .MethodName }}(ctx context.Context, i *pb.{{ .Input.Name }}) (*emptypb.Empty, error) {
{{ template "logger" . }}

{{ template "input" . }}

	logger.Trace().Msg("create reqCtx")
	reqCtx := context.WithValue(ctx, "logger", logger)

	logger.Trace().Msg("call usecase")
	err = self.{{ .DomainCamel }}Usecase.{{ .MethodName }}(reqCtx, uscI)
	if err != nil {
		logger.Warn().Err(err).Msg("usecase err")
		return nil, err
	}

	logger.Trace().Msg("end")
	return &emptypb.Empty{}, nil
}
{{ end -}}

{{- define "method-with-output" }}
func (self *{{ .DomainCamel }}Controller) {{ .MethodName }}(ctx context.Context, _ *emptypb.Empty) (*pb.{{ .Output.Name }}, error) {
{{ template "logger" . }}

	logger.Trace().Msg("create reqCtx")
	reqCtx := context.WithValue(ctx, "logger", logger)

	logger.Trace().Msg("call usecase")
	result, err := self.{{ .DomainCamel }}Usecase.{{ .MethodName }}(reqCtx)
	if err != nil {
		logger.Warn().Err(err).Msg("usecase err")
		return nil, err
	}

{{ template "output" . }}
}
{{ end -}}

{{- define "method-with-input-and-output" }}
func (self *{{ .DomainCamel }}Controller) {{ .MethodName }}(ctx context.Context, i *pb.{{ .Input.Name }}) (*pb.{{ .Output.Name }}, error) {
{{ template "logger" . }}

{{ template "input" . }}

	logger.Trace().Msg("create reqCtx")
	reqCtx := context.WithValue(ctx, "logger", logger)

	logger.Trace().Msg("call usecase")
	result, err := self.{{ .DomainCamel }}Usecase.{{ .MethodName }}(reqCtx, uscI)
	if err != nil {
		logger.Warn().Err(err).Msg("usecase err")
		return nil, err
	}

{{ template "output" . }}
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

	logger   {{ .SpacingRelativeToDomainName }}zerolog.Logger
	validator{{ .SpacingRelativeToDomainName }}adapters.Validator
	{{ .DomainCamel }}Usecase  {{ .DomainSnake }}_usecase.{{ .Domain }}Usecase
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
	Server   {{ .SpacingRelativeToDomainName }}*grpc.Server
	Logger   {{ .SpacingRelativeToDomainName }}zerolog.Logger
	Validator{{ .SpacingRelativeToDomainName }}adapters.Validator
	{{ .Domain }}Usecase  {{ .DomainSnake }}_usecase.{{ .Domain }}Usecase
}

func Add{{ .Domain }}Controller(i *Add{{ .Domain }}ControllerInput) {
	pb.Register{{ .Domain }}Server(i.Server, &{{ .DomainCamel }}Controller{
		logger:   {{ .SpacingRelativeToDomainName }}i.Logger,
		validator:{{ .SpacingRelativeToDomainName }}i.Validator,
		{{ .DomainCamel }}Usecase: i.{{ .Domain }}Usecase,
	})
}
`