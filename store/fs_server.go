package store

import (
	"io/fs"
	"net/http"
	"path/filepath"
)

type ctxKey string

const StoreDirKey ctxKey = "storeDir"

// FSServer returns a fileserver from the store
func FSServer(s fs.FS) http.Handler {
	h := http.FileServer(http.FS(noDirFS{s}))
	return &fsServer{h}
}

type fsServer struct {
	handler http.Handler
}

func (s *fsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodHead:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Add directory prefix if necessary
	storeDir, ok := r.Context().Value(StoreDirKey).(string)
	if ok && storeDir != "" {
		r.URL.Path = filepath.Join(storeDir, r.URL.Path)
	}

	s.handler.ServeHTTP(w, r)
}

type noDirFS struct {
	wrapped fs.FS
}

func (n noDirFS) Open(name string) (fs.File, error) {
	f, err := n.wrapped.Open(name)
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, fs.ErrNotExist
	}

	return f, nil
}
