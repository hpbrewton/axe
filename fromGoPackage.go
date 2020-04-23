package main

import (
	"fmt"
	"path/filepath"
	"os"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
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

// take a file and a comment map, and add the file's comments to the map
func addDeclsFromFile(commentMap map[token.Pos]*ast.CommentGroup, parsedFile *ast.File) {
	// add declaration positions to map
	for _, decl := range parsedFile.Decls {
		switch decl.(type) {
		case *ast.GenDecl:
			gendecl := decl.(*ast.GenDecl)
			if gendecl.Tok == token.IMPORT {
				continue
			}
			for _, spec := range gendecl.Specs{
				if typ, ok := spec.(*ast.TypeSpec); ok {
					commentMap[typ.Name.NamePos] = gendecl.Doc
				}
				if vspec, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range vspec.Names {
						commentMap[name.NamePos] = gendecl.Doc // TODO this is not quite right because the comments get weird with comments, come back 
					}
				}
			}
		case *ast.FuncDecl: 
			funcDecl := decl.(*ast.FuncDecl)
			commentMap[funcDecl.Name.NamePos] = funcDecl.Doc
		}
	}
}

// takes files and returns a list of pathes
func getTypes(pathes []string) ([]*GoFragment, error) {
	fset := token.NewFileSet()

	astFiles := make([]*ast.File, len(pathes))
	commentMap := make(map[token.Pos]*ast.CommentGroup)
	count := 0
	for _, path := range pathes {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}

		// read in text of file
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		// add file to set and ast list
		parsedFile, err := parser.ParseFile(fset, path, contents, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		astFiles[count] = parsedFile
		count++

		addDeclsFromFile(commentMap, parsedFile)
	}
	astFiles = astFiles[:count]

	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check("bob?", fset, astFiles, nil)
	if err != nil {
		return nil, err
	}

	pkgScope := pkg.Scope()

	// now to gather fragments... 
	conv := &GoToAxeConverter {
		Named: make(map[*types.TypeName]Type),
	}
	fragments := make([]*GoFragment, len(pkgScope.Names()))
	for i, name := range pkgScope.Names() {
		var fragment GoFragment
		object := pkgScope.Lookup(name)

		// getting comments
		commentGroup, ok := commentMap[object.Pos()]
		if ok {
			fragment.Comment = commentGroup.Text()
		}

		fragment.pkg = pkg

		// getting file location
		fragment.url = fmt.Sprintf("%v/%v", pkg.Name(), object.Name())

		// getting type
		fragment.Typ = conv.GoTypeToAxeType(object.Type())

		// TODO get position // Or maybe not, just give name

		fragments[i] = &fragment
	}

	return fragments, nil
}

// infrastructure for loading from Go package
func GoFragmentsFromDirectory(root string) ([]*GoFragment, error) {
	fragments := make([]*GoFragment, 0)

	err := filepath.Walk(root, func (path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		// for now,
		// let's ignore these directoreis they do wierd things where they call to C
		// notice /net and /os are in here. Yikes! this needs to change TODO
		if strings.HasPrefix(path, "data/go/src/builtin") ||
		strings.HasPrefix(path, "data/go/src/bytes") ||
		strings.HasPrefix(path, "data/go/src/cmd") || 
		strings.HasPrefix(path, "data/go/src/crypto/ecdsa") ||
		strings.HasPrefix(path, "data/go/src/crypto/tls") ||
		strings.HasPrefix(path, "data/go/src/crypto/x509") || 
		strings.HasPrefix(path, "data/go/src/database/sql") || 
		strings.HasPrefix(path, "data/go/src/debug") || 
		strings.HasPrefix(path, "data/go/src/go/build") || 
		strings.HasPrefix(path, "data/go/src/go/internal") ||
		strings.HasPrefix(path, "data/go/src/net") ||
		strings.HasPrefix(path, "data/go/src/os") ||
		strings.HasPrefix(path, "data/go/src/plugin") ||
		strings.HasPrefix(path, "data/go/src/runtime") ||
		strings.HasPrefix(path, "data/go/src/strings") ||
		strings.HasPrefix(path, "data/go/src/vendor") ||
		strings.Contains(path, "testdata") {
			return nil
		}

		files, err := getGoFiles(path)
		if err != nil {
			return nil //not a go directory, no worry
		}
		frags, err := getTypes(files)
		if err != nil {
			return err 
		}
		fragments = append(fragments, frags...)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return fragments, nil
}

