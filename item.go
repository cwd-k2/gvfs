package gvfs

type ItemKind int

const (
	FileItem      ItemKind = iota
	DirectoryItem ItemKind = iota
)

type Item interface {
	// Commit itself to the real file system.
	Commit(string) error
	Kind() ItemKind
	Name() string
}
