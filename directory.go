package gvfs

import (
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
		Contents: []Item{},
	}
}

// Create an entity under the specified directory
func (d *Directory) Commit(parent string) error {
	dirname := filepath.Join(parent, d.BaseName)
	// mkdir -p in case of the parent directory doesn't exist
	if err := os.MkdirAll(dirname, os.ModePerm); err != nil {
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

// Deprecated: AttachFile is now named CreateFile
func (d *Directory) AttachFile(path *Path) (*File, error) {
	return d.CreateFile(path)
}

// Appending the Item(s) as the Directory's contents.
// Files will be overwritten.
// Directories will be merged.
// File and Directory that have the same path, untouched.
// TODO: Test
func (d *Directory) AppendItem(items ...Item) {
item_loop:
	for _, item := range items {
		for _, content := range d.Contents {
			// check for same name
			// dir  & dir  -> recurse AppendItem
			// file & file -> overwrite
			// otherwise   -> ignore
			if content.Name() == item.Name() {
				switch item := item.(type) {
				case *Directory: // recurse
					if content, ok := content.(*Directory); ok {
						content.AppendItem(item.Contents...)
					}
				case *File: // override
					if _, ok := content.(*File); ok {
						content = item
					}
				}
				// skip anyway
				continue item_loop
			}
		}
		// you can safely append this item here
		d.Contents = append(d.Contents, item)
	}
}

// Create a new File at the given path to this Directory
// returns the attached File
func (d *Directory) CreateFile(path *Path) (*File, error) {
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
			d.Contents = append(d.Contents, item)
		} else {
			item = NewFile(path.Head)
			d.Contents = append(d.Contents, item)
		}
	}

	if path.Next == nil {
		if file, ok := item.(*File); ok {
			return file, nil
		} else {
			return nil, fmt.Errorf("Cannot create a File %s. A Directory named %s already exists.", path.Identity, path.Identity)
		}
	}

	if subdir, ok := item.(*Directory); ok {
		return subdir.CreateFile(path.Next) // go recurse
	} else {
		return nil, fmt.Errorf("Cannot create a Directory %s. A File named %s already exists.", path.Identity, path.Identity)
	}

}

// Create a new Directory at the given path to this Directory
// returns the created Directory
func (d *Directory) CreateDirectory(path *Path) (*Directory, error) {
	var item Item = nil

	for _, c := range d.Contents {
		if c.Name() == path.Head {
			item = c
		}
	}

	// item shouldn't be nil after this part
	if item == nil {
		item = NewDirectory(path.Head)
		d.Contents = append(d.Contents, item)
	}

	if path.Next == nil {
		if dir, ok := item.(*Directory); ok {
			return dir, nil
		} else {
			return nil, fmt.Errorf("Cannot create a Directory %s. A File named %s already exists.", path.Identity, path.Identity)
		}
	}

	if subdir, ok := item.(*Directory); ok {
		return subdir.CreateDirectory(path.Next) // go recurse
	} else {
		return nil, fmt.Errorf("Cannot creat a Directory %s. A File named %s already exists.", path.Identity, path.Identity)
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
		return nil, fmt.Errorf("Item %s not found.", path.Identity)
	}

	if subdir, ok := item.(*Directory); ok {
		return subdir.Search(path.Next) // go recurse
	} else {
		return nil, fmt.Errorf("Item %s is not a Directory. Cannot go deeper.", path.Identity)
	}
}

// Search a File that has the given path
func (d *Directory) SearchFile(path *Path) (*File, error) {
	item, err := d.Search(path)
	if err != nil {
		return nil, err
	}

	if file, ok := item.(*File); ok {
		return file, nil
	} else {
		return nil, fmt.Errorf("%s is not a File, but a Directory.", path.Identity)
	}
}

// Search a Directory that has the given path
func (d *Directory) SearchDirectory(path *Path) (*Directory, error) {
	item, err := d.Search(path)
	if err != nil {
		return nil, err
	}

	if directory, ok := item.(*Directory); ok {
		return directory, nil
	} else {
		return nil, fmt.Errorf("%s is not a Directory, but a File.", path.Identity)
	}
}
