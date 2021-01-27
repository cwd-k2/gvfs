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

	src := gvfs.NewRoot(srcDir)
	dst := gvfs.NewRoot(dstDir)

	pat := regexp.MustCompile(os.Args[3])

	e, err := src.ToItem(pat)
	if err != nil {
		panic(err)
	}

	switch d := e.(type) {
	case *gvfs.Directory:
		for _, e := range d.Contents {
			if err := dst.WriteItem(e); err != nil {
				println(err.Error())
			}
		}
	}

}
