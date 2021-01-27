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

func (f *File) Commit(parent string) error {
	fp, err := os.Create(filepath.Join(parent, f.BaseName))
	if err != nil {
		return err
	}
	defer fp.Close()

	b := bytes.NewBuffer(f.Contents)
	if _, err := b.WriteTo(fp); err != nil {
		return nil
	}

	return nil
}

func (f *File) Kind() ItemKind {
	return FileItem
}
