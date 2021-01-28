package gvfs

import (
	"bytes"
	"os"
	"path/filepath"
)

type File struct {
	BaseName string
	Contents []byte
}

func NewFile(basename string) *File {
	return &File{
		BaseName: basename,
		Contents: []byte{},
	}
}

// Create an entity under the specified directory.
func (f *File) Commit(parent string) error {
	fp, err := os.Create(filepath.Join(parent, f.BaseName))
	if err != nil {
		return err
	}
	defer fp.Close()

	b := bytes.NewBuffer(f.Contents)
	if _, err := b.WriteTo(fp); err != nil {
		return err
	}

	return nil
}

func (f *File) Kind() ItemKind {
	return FileItem
}

func (f *File) Name() string {
	return f.BaseName
}
