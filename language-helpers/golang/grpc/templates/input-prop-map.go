package templates

import "strings"

type InputPropMapTemplProp struct {
	Name    string
	Spacing string
	Value   string
}

type InputPropMapTemplInput struct {
	IndentationLvl int // internal use, controls sub levels

	HasOutput            bool
	MethodName           string
	Optional             bool
	OriginalVariableName string
	Prepare              []string
	Props                []*InputPropMapTemplProp
	Type                 string
	TypePkg              *string
	VarName              string
}

func (self *InputPropMapTemplInput) Idt() string {
	return strings.Repeat("	", self.IndentationLvl)
}

const InputPropMapTempl = `	var {{ .VarName }} *{{ if .TypePkg }}{{ .TypePkg }}.{{ end }}{{ .Type }} = nil
	if {{ .OriginalVariableName }} != nil {
		{{- if .Prepare }}
			{{- range .Prepare }}
{{ . }}
			{{- end }}
		{{- end }}
		{{ .VarName }} = &{{ if .TypePkg }}{{ .TypePkg }}.{{ end }}{{ .Type }}{
			{{- range .Props }}
		 	{{ .Name }}:{{ .Spacing }} {{ .Value }},
			{{- end }}
		}
	}{{ if not .Optional }} else {
		return {{ if .HasOutput }}nil, {{ end }}errors.New("\"{{ .MethodName }}\": \"{{ .OriginalVariableName }}\" must not be nil")
	}{{ end }}`
