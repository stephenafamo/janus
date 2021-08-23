package source

import (
	"fmt"
	"io/fs"
	"strings"
)

// AferoTemplates is an implementation of the Templates interface
// based on aferto
type FsTemplates struct {
	FS         fs.FS
	Suffix     string // required suffix for template files
	TrimSuffix bool   // should suffix be removed from template names?
}

// Walk imiplements the Templates interface
func (p FsTemplates) Walk(walkFunc func(string, fs.File) error) error {
	return fs.WalkDir(p.FS, ".", func(path string, info fs.DirEntry, err error) error {

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

		file, err2 := p.FS.Open(path)
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
