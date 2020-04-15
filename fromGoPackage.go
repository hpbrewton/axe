package main

import (
	"os"
	"go/ast"
	"go/build"
	"go/token"
	"go/importer"
	"go/types"
	"go/parser"
	"strings"
	"log"
	"io/ioutil"
	"path/filepath"
)

// walks from directory and gets all go files
func getGoFiles(root string) ([]string, error) {
	files := make([]string, 0)

	err := filepath.Walk(root, func(path string, _ os.FileInfo, e error) error{
		if e != nil {
			return e
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err 
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
	log.Println(build.Default)
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