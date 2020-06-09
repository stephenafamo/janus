package source

import (
	"net/http"
)

// Templates is an interface containing the raw templates
type Templates interface {
	Walk(func(path string, file http.File) error) error
}
