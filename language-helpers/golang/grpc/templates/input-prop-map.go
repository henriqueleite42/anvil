package templates

import "strings"

type InputPropMapTemplProp struct {
	Name    string
	Spacing string
	Value   string
}

type InputPropMapTemplInput struct {
	IndentationLvl int // internal use, controls sub levels

	MethodName           string
	TypePkg              *string
	OriginalVariableName string
	Prepare              []string
	VarName              string
	Type                 string
	Props                []*InputPropMapTemplProp
	Optional             bool
	HasOutput            bool
}

func (self *InputPropMapTemplInput) Idt() string {
	return strings.Repeat("	", self.IndentationLvl)
}

const InputPropMapTempl = `	var {{ .VarName }} *{{ if ne .TypePkg "" }}{{ .TypePkg }}.{{ end }}{{ .Type }} = nil
	if {{ .OriginalVariableName }} != nil {
							{{ if .Prepare -}}
								{{ range .Prepare -}}
									{{- . }}
								{{- end }}
							{{- end }}
		{{ .VarName }} = &{{ if ne .TypePkg "" }}{{ .TypePkg }}.{{ end }}{{ .Type }}{
			{{- range .Props }}
		 	{{ .Name }}:{{ .Spacing }} {{ .Value }},
			{{- end }}
		}
	}{{ if not .Optional }} else {
		return {{ if .HasOutput }}nil, {{ end }}errors.New("\"{{ .MethodName }}\": \"{{ .OriginalVariableName }}\" must not be nil")
	}{{ end }}`
