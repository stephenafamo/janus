package executor

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"

	"github.com/stephenafamo/janus/views/source"
)

// Executor is an interface that render templates
type Executor interface {
	// Render will execute a template with the data provided
	Render(wr io.Writer, name string, data interface{}) error

	// Exists checks for the presence of a template
	Exists(name string) bool

	// Add adds a new template to the executor
	Add(name string, data io.Reader) error
}

// helper function to load templates into an executor
func loadTemplates(tpl Executor, t source.Templates) error {
	return t.Walk(func(path string, file fs.File) error {
		if file == nil {
			return nil
		}

		finfo, err := file.Stat()
		if err != nil {
			return fmt.Errorf("error getting file info: %w", err)
		}

		if finfo.IsDir() {
			return nil
		}

		err = tpl.Add(path, file)
		if err != nil {
			return fmt.Errorf("error adding new template: %w", err)
		}
		return nil
	})
}

func newTemplate(addFuncs template.FuncMap) *template.Template {
	return template.New("Views").Funcs(funcMap).Funcs(addFuncs)
}
