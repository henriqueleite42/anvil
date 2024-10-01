package templates

type InputPropMapTemplProp struct {
	Name    string
	Spacing string
	Value   string
}

type InputPropMapTemplInput struct {
	MethodName           string
	TypePkg              string
	OriginalVariableName string
	VarName              string
	Type                 string
	Props                []*InputPropMapTemplProp
	Optional             bool
	HasOutput            bool
}

const InputPropMapTempl = `	var {{ .VarName }} *{{ if ne .TypePkg "" }}{{ .TypePkg }}.{{ end }}{{ .Type }} = nil
	if {{ .OriginalVariableName }} != nil {
		{{ .VarName }} = &{{ if ne .TypePkg "" }}{{ .TypePkg }}.{{ end }}{{ .Type }}{
			{{- range .Props }}
		 	{{ .Name }}:{{ .Spacing }} {{ .Value }},
			{{- end }}
		}
	}{{ if not .Optional }} else {
		return {{ if .HasOutput }}nil, {{ end }}errors.New("\"{{ .MethodName }}\": \"{{ .OriginalVariableName }}\" must not be nil")
	}{{ end }}`
