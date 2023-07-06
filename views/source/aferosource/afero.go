package source

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/spf13/afero"
)

// Templates is an implementation of the Templates interface
// based on afero
type Templates struct {
	Store      afero.Fs
	Suffix     string // required suffix for template files
	TrimSuffix bool   // should suffix be removed from template names?
}

// Walk imiplements the Templates interface
func (p Templates) Walk(walkFunc func(string, fs.File) error) error {
	return afero.Walk(p.Store, ".", func(path string, info os.FileInfo, err error) error {
		// Ignore hidden files (files that start with a period)
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// Skip files that do not have the required suffix
		if !strings.HasSuffix(path, p.Suffix) {
			return nil
		}

		if err != nil {
			return err
		}

		file, err2 := p.Store.Open(path)
		if err2 != nil {
			return err2
		}
		defer file.Close()

		if p.TrimSuffix {
			path = strings.TrimSuffix(path, p.Suffix)
		}

		err = walkFunc(path, file)
		if err != nil {
			return fmt.Errorf("error in walkFunc: %w", err)
		}

		return nil
	})
}
