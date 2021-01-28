package main

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/cwd-k2/gvfs"
)

func main() {
	if len(os.Args) != 4 {
		println("usage: filter-cp <srcdir> <dstdir> <ignore-pattern>")
		return
	}

	srcDir, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}

	dstDir, err := filepath.Abs(os.Args[2])
	if err != nil {
		panic(err)
	}

	pat := regexp.MustCompile(os.Args[3])
	src, err := gvfs.Traverse(srcDir, pat)
	if err != nil {
		panic(err)
	}

	for _, content := range src.Contents {
		if err := content.Commit(dstDir); err != nil {
			println(err)
		}
	}
}
