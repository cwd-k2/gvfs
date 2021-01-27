package gvfs

import (
	"os"
	"path/filepath"
)

type Directory struct {
	BaseName string
	Contents []Item
}

// Create an entity under the specified directory
func (d *Directory) Commit(parent string) error {
	dirname := filepath.Join(parent, d.BaseName)

	if err := os.MkdirAll(dirname, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}

	for _, contents := range d.Contents {
		err := contents.Commit(dirname)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Directory) Kind() ItemKind {
	return DirectoryItem
}
