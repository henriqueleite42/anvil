package template

import (
	"fmt"
	"text/template"
)

type templateManager struct {
	templates map[string]*template.Template
}

type localTmpl struct {
	content string
}

func (self *localTmpl) Write(p []byte) (int, error) {
	self.content += string(p)
	return len(p), nil
}

func (self *templateManager) AddTemplate(name string, tmplData string) error {
	if _, ok := self.templates[name]; ok {
		return fmt.Errorf("template \"%s\" already registered", name)
	}

	templ, err := template.New(name).ParseGlob(tmplData)
	if err != nil {
		return err
	}

	self.templates[name] = templ

	return nil
}

func (self *templateManager) Parse(name string, placeholders any) (string, error) {
	templ, ok := self.templates[name]
	if !ok {
		return "", fmt.Errorf("template \"%s\" notfound", name)
	}

	lTempl := &localTmpl{}

	err := templ.Execute(lTempl, placeholders)
	if err != nil {
		return "", err
	}

	return lTempl.content, nil
}

func NewTemplateManager() TemplateManager {
	return &templateManager{
		templates: map[string]*template.Template{},
	}
}
