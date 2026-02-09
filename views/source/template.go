package source

import "io/fs"

// Templates is an interface containing the raw templates
type Templates interface {
	Walk(func(path string, file fs.File) error) error
}
