package gvfs

import (
	"path/filepath"
	"strings"
)

type Path struct {
	Head string
	Next *Path
}

func NewPath(path string) *Path {
	if path == "" {
		return &Path{Head: "", Next: nil}
	}

	list := strings.Split(filepath.Clean(path), string(filepath.Separator))
	return newpath(list)
}

func newpath(ps []string) *Path {
	if len(ps) > 1 {
		return &Path{Head: ps[0], Next: newpath(ps[1:])}
	} else {
		return &Path{Head: ps[0], Next: nil}
	}
}
