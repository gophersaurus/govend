package main

//go:generate go run mkpkgs.go

import (
	"go/parser"
	"go/token"
	"log"
	"strconv"
	"strings"

	"github.com/kr/fs"
)

// scan walks a golang project directory, collecting external package imports.
func scan(dir string) ([]string, error) {

	// define an empty package list to fill
	var pkglist []string

	// create a new walk
	w := fs.Walk(dir)

	// start the walk down the directory tree
	for w.Step() {

		// determine the file statistics once
		fstat := w.Stat()

		// check if we currently are at a directory
		if fstat.IsDir() {

			// check if that directory is "_vendor"
			if fstat.Name() == vendorDir {

				// skip the directory
				w.SkipDir()
				continue
			}
			continue
		}

		// check for errors
		if w.Err() != nil {
			log.Println("govend scan:", w.Err())
			continue
		}

		// check the file is a .go file
		if strings.HasSuffix(w.Path(), ".go") {

			// create an empty fileset
			fset := token.NewFileSet()

			// parse only the import declarations in the .go file
			f, err := parser.ParseFile(fset, w.Path(), nil, parser.ImportsOnly)
			if err != nil {

				// define empty .go file message
				msg := "expected 'package', found 'EOF'"

				// get the error as a string
				e := err.Error()

				// ensure we don't run into memory length issues
				if len(e) >= len(msg) {

					// check for empty .go fiel message
					if e[len(e)-len(msg):] == msg {
						continue
					}
				}
				return nil, err
			}

			// iterate through the files import paths
			for _, i := range f.Imports {

				// unquote the import path value
				importpath, err := strconv.Unquote(i.Path.Value)
				if err != nil {
					return nil, err
				}

				// iterate through the known external packages
				for _, pkg := range pkglist {

					// check if package path already exists, skip the append
					if importpath == pkg {
						goto SKIP
					}
				}

				// if the import path doens't exists in pkgs, add it
				pkglist = append(pkglist, importpath)

			SKIP: // skips the appending of packages that are already present
			}
		}
	}

	return pkglist, nil
}