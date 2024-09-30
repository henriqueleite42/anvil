package templates

type InputPropMapTemplProp struct {
	Name    string
	Spacing string
	Value   string
}

type InputPropMapTemplInput struct {
	VarName string
	Type    string
	Props   []*InputPropMapTemplProp
}

const InputPropMapTempl = `	{{ .VarName }} := &pb.{{ .Type }}{
	 	{{- range .Props }}
		{{ .Name }}:{{ .Spacing }} {{ .Value }},
	 	{{- end }}
	}`
