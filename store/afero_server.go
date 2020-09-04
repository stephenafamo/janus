package store

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

type ctxKey string

const StoreDirKey ctxKey = "storeDir"

// AferoServer returns a fileserver from the store
func AferoServer(s afero.Fs) http.Handler {
	return &aferoServer{s}
}

type aferoServer struct {
	Store afero.Fs
}

func (s *aferoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	errCode := http.StatusInternalServerError

	path := strings.TrimPrefix(r.URL.Path, "/") // remove leading slash if present

	// Add directory prefix if necessary
	storeDir, ok := r.Context().Value(StoreDirKey).(string)
	if ok && storeDir != "" {
		path = filepath.Join(storeDir, path)
	}

	file, err := s.Store.Open(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("ERROR: Unable to get file %v", err)
		http.Error(w, http.StatusText(errCode), errCode)
		return
	}
	if errors.Is(err, os.ErrNotExist) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	info, err := file.Stat()
	if err != nil {
		log.Printf("ERROR: Unable to get file info %v", err)
		http.Error(w, http.StatusText(errCode), errCode)
		return
	}
	if info.IsDir() {
		log.Printf("INFO: attempting to read storage directory")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
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
	} else if len(ctypes) > 0 {
		ctype = ctypes[0]
	}

	w.Header().Set("Content-Type", ctype)
	w.WriteHeader(http.StatusOK)

	if r.Method != "HEAD" {
		_, err = w.Write(fileBytes)
		if err != nil {
			log.Printf("ERROR: problems writing file content to http response writer: %v", err)
		}
	}
}
