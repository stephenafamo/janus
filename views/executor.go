package views

import (
	"html/template"
	"io"
)

// TemplateExecutor is an interface that render templates
type TemplateExecutor interface {
	// Render will execute a template with the data provided
	Render(wr io.Writer, name string, data interface{}) error

	// Lookup checks for the presence of a template
	// Returns nil if it does not exist
	Lookup(name string) *template.Template

	// New adds a new template
	New(name string) *template.Template
}

// NewProdExecutor returns *ProdTemplateExecutor
// it implements the Executor interface
func NewProdExecutor(tpls templates, funcs template.FuncMap) (*ProdTemplateExecutor, error) {

	t := &ProdTemplateExecutor{newTemplate(funcs)}

	err := loadTemplates(t, tpls)
	if err != nil {
		return t, err
	}

	return t, nil
}

// NewHotExecutor returns a template executor
// that implements the Executor interface
// If will parse all the templates before rendering each time.
// This is slower but useful during development
func NewHotExecutor(tpls templates, funcs template.FuncMap) (*HotTemplateExecutor, error) {

	t := &HotTemplateExecutor{Template: newTemplate(funcs), files: tpls}

	err := loadTemplates(t, tpls)
	if err != nil {
		return t, err
	}

	return t, nil
}

// ProdTemplateExecutor is a Template Executor optimised for production
// It only loads the templates once. The program has to restart to reload the templates
type ProdTemplateExecutor struct {
	*template.Template
}

// Render implements the templateExecutor interface
func (p ProdTemplateExecutor) Render(wr io.Writer, name string, data interface{}) error {
	return p.ExecuteTemplate(wr, name, data)
}

// HotTemplateExecutor is a Template Executor meant for development
// It will reload the templates before each render
type HotTemplateExecutor struct {
	*template.Template
	files templates
}

// Render implements the templateExecutor interface
func (h *HotTemplateExecutor) Render(wr io.Writer, name string, data interface{}) error {
	err := loadTemplates(h, h.files)
	if err != nil {
		return err
	}

	d, err := h.Clone()
	if err != nil {
		return err
	}

	return d.ExecuteTemplate(wr, name, data)
}

func newTemplate(addFuncs template.FuncMap) *template.Template {
	return template.New("Views").Funcs(funcMap).Funcs(addFuncs)
}
