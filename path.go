package gvfs

import (
	"path/filepath"
	"strings"
)

type Path struct {
	Identity string
	Head     string
	Next     *Path
}

func NewPath(path string) *Path {
	// strings.Split() returns 1> elements slice
	list := strings.Split(filepath.Clean(path), string(filepath.Separator))
	return newpath("", list)
}

func newpath(parent string, ps []string) *Path {
	// surely, len(ps) > 0
	dirname := filepath.Join(parent, ps[0])
	if len(ps) > 1 {
		return &Path{
			Identity: dirname,
			Head:     ps[0],
			Next:     newpath(dirname, ps[1:]),
		}
	} else {
		return &Path{
			Identity: dirname,
			Head:     ps[0],
			Next:     nil,
		}
	}
}
