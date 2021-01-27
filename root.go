package gvfs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type Root struct {
	Path string
}

func NewRoot(path string) *Root {
	return &Root{Path: path}
}

func (r *Root) ToItem(ignore *regexp.Regexp) (Item, error) {
	// ensure the r.Path points an existing directory
	if info, err := os.Stat(r.Path); err != nil {
		return nil, err
	} else if !info.IsDir() {
		return nil, errors.New(fmt.Sprintf("%s is not a directory.", r.Path))
	}

	e, err := traverse(r.Path, ignore)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (r *Root) WriteItem(e Item) error {
	// err if the directory couldn't be created somehow
	if err := os.MkdirAll(r.Path, 0755); err != nil && !os.IsExist(err) {
		return err
	}

	if err := e.Commit(r.Path); err != nil {
		println(err.Error())
	}

	return nil
}

func traverse(path string, ignore *regexp.Regexp) (Item, error) {
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
