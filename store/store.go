package store

import (
	"io"
)

// Store represents an object store
type Store interface {
	AddFile(path string, file io.Reader) error
	UpdateFile(path string, file io.Reader) error
	AddOrUpdateFile(path string, file io.Reader) error
	GetFile(path string) (io.Reader, error)
	DeleteFile(path string) error
}
