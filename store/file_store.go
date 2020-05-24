package store

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileStore is an object store based on a file directory
type FileStore struct {
	Directory string
}

func (f FileStore) createDirIfNotExist(dir string) error {
	_, err := os.Stat(dir)

	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				return err // error when creating directory
			}
			return nil // all went according to plan
		}

		return err // some other error from os.Stat
	}

	return nil // directory exists
}

func (f FileStore) doesFileExist(path string) (bool, error) {
	filename := filepath.Join(f.Directory, path)

	_, err := os.Stat(filename)

	if err != nil && os.IsNotExist(err) {
		return false, nil // files does not exist
	}

	return true, err
}

// AddFile adds a file to the store
func (f FileStore) AddFile(path string, file io.Reader) error {
	var err error

	exists, err := f.doesFileExist(path)

	if err != nil {
		return err
	}
	if exists != false {
		return errors.New("file: '" + path + "' already exists")
	}

	filename := filepath.Join(f.Directory, path)

	f.createDirIfNotExist(filepath.Dir(filename))

	newFile, err := os.Create(filename)
	newFile.Close()
	if err != nil {
		return err
	}
	fmt.Printf("Flie added: %#v \n", filename)

	return f.UpdateFile(path, file)
}

// UpdateFile updates a file in the store
func (f FileStore) UpdateFile(path string, file io.Reader) error {
	var err error

	exists, err := f.doesFileExist(path)

	if err != nil {
		return err
	}
	if exists != true {
		return errors.New("file: '" + path + "' does not exists")
	}

	filename := filepath.Join(f.Directory, path)

	newFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer newFile.Close()
	io.Copy(newFile, file)

	return nil
}

// AddOrUpdateFile adds or updates a file in the store
func (f FileStore) AddOrUpdateFile(path string, file io.Reader) error {
	var err error
	filename := filepath.Join(f.Directory, path)

	f.createDirIfNotExist(filepath.Dir(filename))

	newFile, err := os.Create(filename)
	newFile.Close()
	if err != nil {
		return err
	}

	return f.UpdateFile(path, file)
}

// DeleteFile deletes a file from the store
func (f FileStore) DeleteFile(path string) error {
	return os.Remove(path)
}

// GetFile gets a file from the store
func (f FileStore) GetFile(path string) (io.Reader, error) {
	var file io.Reader
	var err error

	exists, err := f.doesFileExist(path)

	if err != nil {
		return file, err
	}
	if exists != true {
		return file, errors.New("file: '" + path + "' does not exists")
	}

	filename := filepath.Join(f.Directory, path)

	file, err = os.Open(filename)
	if err != nil {
		return file, err
	}

	return file, nil
}
