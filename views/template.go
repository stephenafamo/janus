package views

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

// Templates is an interface containing the templates
type templates interface {
	Walk(func(path string, file http.File) error) error
}

func loadTemplates(tpl TemplateExecutor, t templates) error {

	err := t.Walk(func(path string, file http.File) error {
		if file == nil {
			return nil
		}

		finfo, err := file.Stat()
		if err != nil {
			return err
		}

		if finfo.IsDir() {
			return nil
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}

		nt := tpl.New(
			strings.TrimSuffix(
				path,
				filepath.Ext(finfo.Name()),
			),
		)
		_, err = nt.Parse(string(data))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
