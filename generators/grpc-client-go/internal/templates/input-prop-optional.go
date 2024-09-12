package templates

type InputPropOptionalTemplInput struct {
	VarName              string
	OriginalVariableName string
	Type                 string
	ValueToAssign        string
	NeedsPointer         bool
}

const InputPropOptionalTempl = `	var {{ .VarName }} {{ .Type }} = nil
	if {{ .OriginalVariableName }} != nil {
		{{- if .NeedsPointer }}
		pointer := {{ .ValueToAssign }}
		{{ .VarName }} = &pointer
		{{- else }}
		{{ .VarName }} = {{ .ValueToAssign }}
		{{- end }}
	}`
