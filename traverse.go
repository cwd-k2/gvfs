package gvfs

import (
	"errors"
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
		return nil, errors.New(fmt.Sprintf("%s is not a directory.", path))
	}

	d, err := traverse(path, ignore)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func traverse(path string, ignore *regexp.Regexp) (*Directory, error) {
	matched, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return nil, err
	}

	contents := make([]Item, 0)

	for _, child := range matched {
		if ignore != nil && ignore.MatchString(child) {
			continue
		}

		info, err := os.Stat(child)
		if err != nil {
			println(err.Error())
			continue
		}

		if info.Mode().IsDir() {
			d, err := traverse(child, ignore)

			if err != nil {
				println(err.Error())
				continue
			}

			contents = append(contents, d)
		} else if info.Mode().IsRegular() {
			fp, err := os.Open(child)
			if err != nil {
				println(err.Error())
				continue
			}
			defer fp.Close()

			b, err := ioutil.ReadFile(child)
			if err != nil {
				println(err.Error())
				continue
			}
			f := &File{BaseName: filepath.Base(child), Contents: b}

			contents = append(contents, f)
		}
	}

	return &Directory{BaseName: filepath.Base(path), Contents: contents}, nil
}
