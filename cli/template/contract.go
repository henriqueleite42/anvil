package template

type TemplateManager interface {
	AddTemplate(name string, template string) error
	Parse(name string, data any) (string, error)
}
