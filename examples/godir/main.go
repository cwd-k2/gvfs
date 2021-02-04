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

	target, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}

	dirname, basename := filepath.Split(target)

	newfiles := []string{
		"README.md",
		fmt.Sprintf("cmd/%s/main.go", basename),
	}

	newdires := []string{
		"pkg",
		"test",
		"examples",
		"internal/pkg",
	}

	directory := gvfs.NewDirectory(basename)

	for _, filename := range newfiles {
		if _, err := directory.CreateFile(gvfs.NewPath(filename)); err != nil {
			println("errors: create-file")
			println(err.Error())
		}
	}

	for _, direname := range newdires {
		if _, err := directory.CreateDirectory(gvfs.NewPath(direname)); err != nil {
			println("errors: create-directory")
			println(err.Error())
		}
	}

	if b, err := json.MarshalIndent(directory, "", "  "); err != nil {
		println("errors: json-marshal")
	} else {
		println(string(b))
	}

	if err := directory.Commit(dirname); err != nil {
		panic(err)
	}
}
