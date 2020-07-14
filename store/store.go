package store

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

// Store represents an object store
type Store interface {
	FileExists(path string) bool
	AddFile(path string, file io.Reader) error
	UpdateFile(path string, file io.Reader) error
	AddOrUpdateFile(path string, file io.Reader) error
	GetFile(path string) (io.Reader, error)
	DeleteFile(path string) error
}

// FileServer returns a fileserver from the store
func FileServer(s Store) http.Handler {
	return &storeServer{s}
}

type storeServer struct {
	Store Store
}

func (s *storeServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	errCode := http.StatusInternalServerError

	path := strings.TrimPrefix(r.URL.Path, "/") // remove leading slash if present
	file, err := s.Store.GetFile(path)
	if err != nil {
		log.Printf("ERROR: Unable to get file %v", err)
		http.Error(w, http.StatusText(errCode), errCode)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("ERROR: Unable to read file %v", err)
		http.Error(w, http.StatusText(errCode), errCode)
		return
	}

	// If Content-Type isn't set, use the file's extension to find it, but
	// if the Content-Type is unset explicitly, do not sniff the type.
	ctypes, haveType := w.Header()["Content-Type"]
	var ctype string
	if !haveType {
		ctype = mime.TypeByExtension(filepath.Ext(path))
		if ctype == "" {
			// read a chunk to decide between utf-8 text and binary
			var buf [512]byte
			n, _ := io.ReadFull(bytes.NewBuffer(fileBytes), buf[:])
			ctype = http.DetectContentType(buf[:n])
		}
		w.Header().Set("Content-Type", ctype)
	} else if len(ctypes) > 0 {
		ctype = ctypes[0]
	}

	w.WriteHeader(http.StatusOK)

	if r.Method != "HEAD" {
		w.Write(fileBytes)
	}
}
