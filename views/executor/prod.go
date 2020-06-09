package executor

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"

	"github.com/stephenafamo/janus/views/source"
)

// Prod is a Template Executor optimised for production
// It only loads the templates once. The program has to restart to reload the templates
type Prod struct {
	*template.Template
}

// NewProd returns *ProdTemplateExecutor
// it implements the Executor interface
func NewProd(tpls source.Templates, funcs template.FuncMap) (*Prod, error) {

	t := &Prod{newTemplate(funcs)}

	err := loadTemplates(t, tpls)
	if err != nil {
		return t, err
	}

	return t, nil
}

// Exists checks for the presence of a template
func (p Prod) Exists(name string) bool {
	return p.Lookup(name) != nil
}

// Add adds a new template to the executor
func (p Prod) Add(name string, data io.Reader) error {

	byts, err := ioutil.ReadAll(data)
	if err != nil {
		return fmt.Errorf("error reading bytes from new template: %w", err)
	}

	nt := p.Template.New(name)
	_, err = nt.Parse(string(byts))
	if err != nil {
		return fmt.Errorf("error parsing new template: %w", err)
	}

	return nil
}

// Render implements the templateExecutor interface
func (p Prod) Render(wr io.Writer, name string, data interface{}) error {
	return p.ExecuteTemplate(wr, name, data)
}
