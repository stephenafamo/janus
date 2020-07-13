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
	Root       string
	Suffix     string // required suffix for template files
	TrimSuffix bool   // should suffix be removed from template names?
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

		// Skip files that do not have the required suffix
		if !strings.HasSuffix(path, p.Suffix) {
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

		cleanPath := strings.TrimPrefix(path, dirPrefix)
		if p.TrimSuffix {
			cleanPath = strings.TrimSuffix(cleanPath, p.Suffix)
		}

		return walkFunc(cleanPath, file)
	})
}
