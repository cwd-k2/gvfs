package gvfs

import (
	"os"
	"regexp"
)

// Deprecated: Use gvfs.Traverse or gvfs.NewDirectory.
// Stands for a directory where we can read/write Item.
type Root struct {
	Path string
}

// Deprecated: Use gvfs.Traverse or gvfs.NewDirectory.
func NewRoot(path string) *Root {
	return &Root{Path: path}
}

// Deprecated: Use gvfs.Traverse
// Just convert this Root to Item struct (Directory).
// The returned struct should be a Directory that has the same path to this Root object.
func (r *Root) ToItem(ignore *regexp.Regexp) (*Directory, error) {
	return Traverse(r.Path, ignore)
}

// Deprecated: Use gvfs.Item#Commit directly.
// Write an Item object as a real filesystem entity.
// The object's structure will be written under the Root's path.
func (r *Root) WriteItem(i Item) error {
	// err if the directory couldn't be created somehow
	if err := os.MkdirAll(r.Path, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}

	if err := i.Commit(r.Path); err != nil {
		println(err.Error())
	}

	return nil
}
