package executor

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"sync"

	"github.com/stephenafamo/janus/views/source"
)

// Hot is a Template Executor meant for development
// It will reload the templates before each render
type Hot struct {
	*template.Template
	files source.Templates
	mu    sync.Mutex
}

// NewHot returns a template executor
// that implements the Executor interface
// If will parse all the templates before rendering each time.
// This is slower but useful during development
func NewHot(tpls source.Templates, funcs template.FuncMap) (*Hot, error) {

	t := &Hot{Template: newTemplate(funcs), files: tpls}

	err := t.loadFiles()
	if err != nil {
		return t, err
	}

	return t, nil
}

// Exists checks for the presence of a template
func (h *Hot) Exists(name string) bool {
	return h.Lookup(name) != nil
}

// Add adds a new template to the executor
func (h *Hot) Add(name string, data io.Reader) error {

	byts, err := ioutil.ReadAll(data)
	if err != nil {
		return fmt.Errorf("error reading bytes from new template: %w", err)
	}

	nt := h.Template.New(name)
	_, err = nt.Parse(string(byts))
	if err != nil {
		return fmt.Errorf("error parsing new template: %w", err)
	}

	return nil
}

// Render implements the templateExecutor interface
func (h *Hot) Render(wr io.Writer, name string, data interface{}) error {
	err := h.loadFiles()
	if err != nil {
		return err
	}

	d, err := h.Clone()
	if err != nil {
		return err
	}

	return d.ExecuteTemplate(wr, name, data)
}

func (h *Hot) loadFiles() error {

	// So we do not reload files at the same time
	h.mu.Lock()
	defer h.mu.Unlock()

	err := loadTemplates(h, h.files)
	if err != nil {
		return err
	}

	return nil
}
