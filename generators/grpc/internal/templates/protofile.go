package templates

import "fmt"

type ProtofileTemplInputEnumValue struct {
	Name    string
	Spacing string
	Idx     int32
}

type ProtofileTemplInputEnum struct {
	Name             string
	Values           []*ProtofileTemplInputEnumValue
	DeprecatedValues []int32
}

func (self *ProtofileTemplInputEnum) GetDeprecatedValues() string {
	var result string
	for _, v := range self.DeprecatedValues {
		if len(result) > 0 {
			result += fmt.Sprintf(", %v", v)
		} else {
			result += fmt.Sprintf("%v", v)
		}
	}

	return result
}

type ProtofileTemplInputMethod struct {
	Name   string
	Input  *string
	Output *string
}

type ProtofileTemplInputTypeProp struct {
	Type     string
	Spacing1 string // Spacing between type and name
	Name     string
	Spacing2 string // Spacing between name and idx
	Idx      int
}

type ProtofileTemplInputType struct {
	Name  string
	Props []*ProtofileTemplInputTypeProp
}

type ProtofileTemplInput struct {
	Domain  string
	Imports []string
	Methods []*ProtofileTemplInputMethod
	Enums   []*ProtofileTemplInputEnum
	Types   []*ProtofileTemplInputType
}

const ProtofileTempl = `syntax = "proto3";

{{ range .Imports -}}
import "{{ . }}";
{{ end }}
service {{ .Domain }}Api {
	{{- range .Methods -}}
		{{- if .Input }}
			{{- if .Output }}
	rpc {{ .Name }}({{ .Input }}) returns ({{ .Output }}) {}
			{{- else }}
	rpc {{ .Name }}({{ .Input }}) {}
			{{- end }}
		{{- else }}
			{{- if .Output }}
	rpc {{ .Name }}() returns ({{ .Output }}) {}
			{{- else }}
	rpc {{ .Name }}() {}
			{{- end }}
		{{- end }}
	{{- end }}
}

{{ range .Enums -}}
enum {{ .Name }} {
	{{- range .Values }}
	{{ .Name }}{{ .Spacing }} = {{ .Idx }};
	{{- end }}
	{{- if .DeprecatedValues }}

	reserved {{ .GetDeprecatedValues }};
	{{- end }}
}
{{ end }}
{{ range .Types -}}
message {{ .Name }} {
	{{- range .Props }}
	{{ .Type }}{{ .Spacing1 }} {{ .Name }}{{ .Spacing2 }} = {{ .Idx }};
	{{- end}}
}
{{ end }}`
