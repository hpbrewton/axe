package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"log"
	"strings"
)

// walks from directory and gets all go files
func getGoFiles(root string) ([]string, error) {
	pkg, err := (build.Default).ImportDir(root, 0)
	if err != nil {
		return nil, err
	}

	// we want to get all normal Go Files and Cgo files out of the directory
	files := append(pkg.GoFiles, pkg.CgoFiles...)

	// and we want to append the root to eacch of these
	for i, path := range files {
		files[i] = fmt.Sprintf("%s/%s", root, path)
	}

	return files, nil
}

// takes files and returns a list of pathes
func getTypes(pathes []string) ([]Fragment, error) {
	fset := token.NewFileSet()

	astFiles := make([]*ast.File, len(pathes))
	count := 0
	for i, path := range pathes {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}

		// read in text of file
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println(i, path)
			return nil, err
		}

		// add file to set and ast list
		parsedFile, err := parser.ParseFile(fset, path, contents, 0)
		if err != nil {
			return nil, err
		}
		astFiles[count] = parsedFile
		count++
	}
	astFiles = astFiles[:count]

	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check("bob?", fset, astFiles, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(pkg.Name())
	log.Println(pkg.Scope())

	return nil, nil
}

// infrastructure for loading from Go package
func main() {
	root := "data/go/src/archive/tar"
	files, err := getGoFiles(root)
	log.Println(files[0])

	if err != nil {
		log.Fatal(err)
	}

	_, err = getTypes(files)
	if err != nil {
		log.Println(err)
	}
}
