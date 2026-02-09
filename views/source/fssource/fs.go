package fssource

import (
	"fmt"
	"io/fs"
	"strings"
)

// Templates is an implementation of the source.Templates interface
// based on fs.FS
type Templates struct {
	FS         fs.FS
	Suffixes   []string // required suffixes for template files
	TrimSuffix bool     // should suffix be removed from template names?
}

// Walk imiplements the Templates interface
func (p Templates) Walk(walkFunc func(string, fs.File) error) error {
	return fs.WalkDir(p.FS, ".", func(path string, info fs.DirEntry, err error) error {
		// Ignore hidden files (files that start with a period)
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		var suffix string
		for k, s := range p.Suffixes {
			if strings.HasSuffix(path, s) {
				suffix = s
				break
			}

			// Skip files that do not have any of the required suffixes
			if k == len(p.Suffixes)-1 {
				return nil
			}
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
			path = strings.TrimSuffix(path, suffix)
		}

		err = walkFunc(path, file)
		if err != nil {
			return fmt.Errorf("error in walkFunc: %w", err)
		}

		return nil
	})
}
