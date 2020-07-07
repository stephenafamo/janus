package source

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// FileTemplates is an implementation of the Templates interface
// based on the excellent packr library.
type FileTemplates struct {
	Root string
}

// NewFile return a new instance of FileTemplates when given a root directory
func NewFile(rootDir string) FileTemplates {
	return FileTemplates{
		Root: rootDir,
	}
}

// Walk imiplements the Templates interface
func (p FileTemplates) Walk(walkFunc func(string, http.File) error) error {
	return filepath.Walk(p.Root, func(path string, info os.FileInfo, err error) error {

		// Ignore hidden files (files that start with a period)
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		if err != nil {
			return err
		}

		file, err2 := os.Open(path)
		if err2 != nil {
			return err2
		}

		dirPrefix := filepath.Clean(p.Root) + "/"

		return walkFunc(strings.TrimPrefix(path, dirPrefix), file)
	})
}
