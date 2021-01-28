package gvfs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Directory struct {
	BaseName string
	Contents []Item
}

func NewDirectory(basename string) *Directory {
	return &Directory{
		BaseName: basename,
		Contents: make([]Item, 0),
	}
}

// Create an entity under the specified directory
func (d *Directory) Commit(parent string) error {
	dirname := filepath.Join(parent, d.BaseName)

	if err := os.MkdirAll(dirname, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}

	for _, c := range d.Contents {
		err := c.Commit(dirname)
		if err != nil {
			println(err.Error())
			continue
		}
	}

	return nil
}

func (d *Directory) Kind() ItemKind {
	return DirectoryItem
}

func (d *Directory) Name() string {
	return d.BaseName
}

// Attach a new File at the given path to this Directory
// returns the attached Item
func (d *Directory) AttachFile(path *Path) (*File, error) {
	var item Item = nil

	for _, c := range d.Contents {
		if c.Name() == path.Head {
			item = c
		}
	}

	// item shouldn't be nil after this part
	if item == nil {
		if path.Next != nil {
			item = NewDirectory(path.Head)
		} else {
			item = NewFile(path.Head)
		}
	}

	if path.Next == nil {
		if file, ok := item.(*File); ok {
			d.Contents = append(d.Contents, file)
			return file, nil
		} else {
			return nil, errors.New(fmt.Sprintf("Cannot attach a File %s. A Directory named %s already exists.", path.Identity, path.Identity))
		}
	}

	if subdir, ok := item.(*Directory); ok {
		d.Contents = append(d.Contents, subdir)
		return subdir.AttachFile(path.Next) // go recurse
	} else {
		return nil, errors.New(fmt.Sprintf("Cannot creat a Directory %s. A File named %s already exists.", path.Identity, path.Identity))
	}

}

// Search an Item that has the given path
func (d *Directory) Search(path *Path) (Item, error) {
	var item Item = nil

	for _, c := range d.Contents {
		if c.Name() == path.Head {
			item = c
		}
	}

	if item != nil && path.Next == nil {
		return item, nil
	}

	if item == nil {
		return nil, errors.New(fmt.Sprintf("Item %s not found.", path.Identity))
	}

	if subdir, ok := item.(*Directory); ok {
		return subdir.Search(path.Next) // go recurse
	} else {
		return nil, errors.New(fmt.Sprintf("Item %s is not a Directory. Cannot go deeper.", path.Identity))
	}
}

// Search an File that has the given path
func (d *Directory) SearchFile(path *Path) (*File, error) {
	item, err := d.Search(path)
	if err != nil {
		return nil, err
	}

	if file, ok := item.(*File); ok {
		return file, nil
	} else {
		return nil, errors.New(fmt.Sprintf("%s is not a File, but a Directory.", path.Identity))
	}
}
