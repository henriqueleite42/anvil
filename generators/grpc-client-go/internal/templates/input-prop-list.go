package templates

type InputPropListTemplInput struct {
	MethodName           string
	VarName              string
	OriginalVariableName string
	Type                 string
	ValueToAppend        string
	Optional             bool
	ChildOptional        bool
	HasOutput            bool
}

const InputPropListTempl = `	var {{ .VarName }} []{{ .Type }} = nil
	if {{ .OriginalVariableName }} != nil {
		{{ .VarName }} = make([]{{ .Type }}, 0, cap({{ .OriginalVariableName }}))
		{{- if .ChildOptional }}
		for _, v := range {{ .OriginalVariableName }} {
			if v == nil {
				continue
			}
			{{ .VarName }} = append({{ .VarName }}, {{ .ValueToAppend }})
		}
		{{- else}}
		for k, v := range {{ .OriginalVariableName }} {
			{{ .VarName }}[k] = {{ .ValueToAppend }}
		}
		{{- end }}
	}{{ if not .Optional }} else {
		return {{ if .HasOutput }}nil, {{ end }}errors.New("\"{{ .MethodName }}\": \"{{ .OriginalVariableName }}\" must not be nil")
	}{{ end }}`
