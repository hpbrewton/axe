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
func getTypes(pathes []string) ([]*Fragment, error) {
	fset := token.NewFileSet()

	astFiles := make([]*ast.File, len(pathes))
	commentMap := make(map[token.Pos]*ast.CommentGroup)
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
		parsedFile, err := parser.ParseFile(fset, path, contents, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		astFiles[count] = parsedFile
		count++

		addDeclsFromFile(commentMap, parsedFile)
	}
	astFiles = astFiles[:count]
	// log.Println(commentMap)

	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check("bob?", fset, astFiles, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(pkg.Name())
	pkgScope := pkg.Scope()
	for _, name := range pkgScope.Names() {
		object := pkgScope.Lookup(name)
		comments, ok := commentMap[object.Pos()] 
		if !ok {
			log.Println(object)
		} else {
			log.Println(comments)
		}
	}

	// now to gather fragments... 
	conv := &GoToAxeConverter {
		Named: make(map[*types.TypeName]Type),
	}
	fragments := make([]*Fragment, len(pkgScope.Names()))
	for i, name := range pkgScope.Names() {
		var fragment Fragment
		object := pkgScope.Lookup(name)

		// getting comments
		commentGroup, ok := commentMap[object.Pos()]
		if ok {
			fragment.comment = commentGroup.Text()
		}

		// getting file location
		fragment.url = fmt.Sprintf("%v/%v", pkg.Name(), object.Name())

		// getting type
		fragment.typ = conv.GoTypeToAxeType(object.Type())

		// TODO get position // Or maybe not, just give name

		fragments[i] = &fragment
	}

	return fragments, nil
}

// infrastructure for loading from Go package
func main() {
	root := "data/go/src/archive/tar"
	files, err := getGoFiles(root)

	if err != nil {
		log.Fatal("here,", err)
	}

	fragments, err := getTypes(files)
	if err != nil {
		log.Println(err)
	}

	for _, fragment := range fragments {
		log.Println(fragment.url, fragment.typ)
		switch fragment.typ.(type) {
		case *MethodHaver:
			for name, method := range fragment.typ.(*MethodHaver).methods {
				log.Println("\t", name, method.String())
			}
		case *Interface:
			for name, method := range fragment.typ.(*Interface).methods {
				log.Println("\t", name, method.String())
			}
		}
	}
}
