package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cwd-k2/gvfs"
)

func main() {
	if len(os.Args) != 2 {
		println("usage: godir <dirname>")
		return
	}

	basename := filepath.Base(os.Args[1])
	newfiles := []string{
		"README.md",
		fmt.Sprintf("cmd/%s/main.go", basename),
		"pkg/.gitkeep",
		"test/.gitkeep",
		"examples/.gitkeep",
		"internal/pkg/.gitkeep",
	}

	directory := gvfs.NewDirectory(basename)

	for _, filename := range newfiles {
		if _, err := directory.CreateFile(gvfs.NewPath(filename)); err != nil {
			println("errors: attachfile")
			println(err.Error())
		}
	}

	if b, err := json.MarshalIndent(directory, "", "  "); err != nil {
		println("errors: json-marshal")
	} else {
		println(string(b))
	}

	target, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}

	if err := directory.Commit(filepath.Dir(target)); err != nil {
		panic(err)
	}
}
