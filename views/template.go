package views

import (
	"fmt"
	"net/http"
)

// Templates is an interface containing the raw templates
type Templates interface {
	Walk(func(path string, file http.File) error) error
}

func loadTemplates(tpl TemplateExecutor, t Templates) error {
	return t.Walk(func(path string, file http.File) error {
		if file == nil {
			return nil
		}

		finfo, err := file.Stat()
		if err != nil {
			return fmt.Errorf("error getting filel info: %w", err)
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
