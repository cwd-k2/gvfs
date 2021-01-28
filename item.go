package gvfs

type ItemKind int

const (
	FileItem      ItemKind = iota
	DirectoryItem ItemKind = iota
)

type Item interface {
	Commit(string) error
	Kind() ItemKind
	Name() string
}
