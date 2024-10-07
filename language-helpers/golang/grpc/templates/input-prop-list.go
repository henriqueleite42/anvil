package templates

import "strings"

type InputPropListTemplInput struct {
	IndentationLvl int // internal use, controls sub levels

	ChildOptional        bool
	HasOutput            bool
	MethodName           string
	Optional             bool
	OriginalVariableName string
	Prepare              []string
	Type                 string
	ValueToAppend        string
	VarName              string
}

func (self *InputPropListTemplInput) Idt() string {
	return strings.Repeat("	", self.IndentationLvl)
}

const InputPropListTempl = `	var {{ .VarName }} []{{ .Type }} = nil
	if {{ .OriginalVariableName }} != nil {
		{{ .VarName }} = make([]{{ .Type }}, 0, len({{ .OriginalVariableName }}))
		for _, v := range {{ .OriginalVariableName }} {
			{{- if .ChildOptional }}
			if v == nil {
				continue
			}
			{{- end }}
			{{- if .Prepare }}
				{{- range .Prepare }}
{{ . }}
				{{- end }}
			{{- end }}
			{{ .VarName }} = append({{ .VarName }}, {{ .ValueToAppend }})
		}
	}{{ if not .Optional }} else {
		return {{ if .HasOutput }}nil, {{ end }}errors.New("\"{{ .MethodName }}\": \"{{ .OriginalVariableName }}\" must not be nil")
	}{{ end }}`
