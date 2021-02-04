package gvfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

// Traverse to get a Directory structure at given path.
// The returned struct should be a Directory that has the same path to this Root object.
func Traverse(path string, ignore *regexp.Regexp) (*Directory, error) {
	// ensure the given path points an existing directory
	if info, err := os.Stat(path); err != nil {
		return nil, err
	} else if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory.", path)
	}

	d, err := traverse(path, ignore)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func traverse(path string, ignore *regexp.Regexp) (*Directory, error) {
	children, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	contents := make([]Item, 0)

	for _, child := range children {
		// FileInfo.Name() returns the base name
		name := filepath.Join(path, child.Name())

		if ignore != nil && ignore.MatchString(name) {
			continue
		}

		// Directory, RegularFile
		// TODO: what about SymLink?
		if child.Mode().IsDir() {
			// recursively traverse
			d, err := traverse(name, ignore)
			if err != nil {
				println(err.Error())
				continue
			}
			contents = append(contents, d)
		} else if child.Mode().IsRegular() {
			// read contents directly
			b, err := ioutil.ReadFile(name)
			if err != nil {
				println(err.Error())
				continue
			}
			f := &File{BaseName: child.Name(), Contents: b}
			contents = append(contents, f)
		}
	}

	return &Directory{BaseName: filepath.Base(path), Contents: contents}, nil
}
