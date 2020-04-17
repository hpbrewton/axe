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
		fragment.typ = object.Type()

		// TODO get position // Or maybe not, just give name

		fragments[i] = &fragment
	}

	return fragments, nil
}

func GoVarToType(va *types.Var) Type {
	return GoTypeToAxeType(va.Type())
}

func GoPointerToAxeKind(ptr *types.Pointer) *Kind {
	underType := GoTypeToAxeType(ptr.Elem())
	return &Kind {
		name: "*",
		arguments: []Type{underType},
	}
}

func GoTupleToSliceOfAxeTypes(tup *types.Tuple) []Type {
	lentup := tup.Len()
	sliceTypes := make([]Type, lentup)
	for i := 0; i < lentup; i++ {
		va := tup.At(i)
		sliceTypes[i] = GoVarToType(va)
	}
	return sliceTypes
}

func GoBasicToAxePrimative(basic *types.Basic) *Primative {
	return &Primative{basic.Name()}
}

func GoSignatureToAxeFunction(sig *types.Signature) *Function {
	return &Function {
		arguments: GoTupleToSliceOfAxeTypes(sig.Params()),
		output: GoTupleToSliceOfAxeTypes(sig.Results()),
	}
}

func GoArrayToAxeArray(array *types.Array) *Array {
	underlying := GoTypeToAxeType(array.Elem())
	return &Array {
		size: int(array.Len()),
		typ: underlying,
	}
}

func GoSliceToAxeKind(slice *types.Slice) *Kind {
	underlying := GoTypeToAxeType(slice.Elem())
	return &Kind{
		name: "slice",
		arguments: []Type{underlying},
	}
}

func GoMapToAxeKind(mp *types.Map) *Kind {
	from := GoTypeToAxeType(mp.Key())
	to := GoTypeToAxeType(mp.Elem())
	return &Kind {
		name: "map",
		arguments: []Type{from, to},
	}
}

func GoChanToAxeKind(ch *types.Chan) *Kind {
	underlying := GoTypeToAxeType(ch.Elem())
	return &Kind {
		name: "chan",
		arguments: []Type{underlying},
	}
}

func GoStructToAxeStruct(strct *types.Struct) *Struct {
	nfields := strct.NumFields()
	fieldNames := make([]string, nfields)
	fields := make([]Type, nfields)
	for i := 0; i < nfields; i++ {
		va := strct.Field(i)
		fieldNames[i] = va.Name()
		fields[i] = GoVarToType(va)
	}
	return &Struct {
		fieldNames: fieldNames,
		fields: fields,
	}
}

// func GoNamedtoAxeName(named *types.Named) Named {
	
// }

// converts go types to actual types
func GoTypeToAxeType(gType types.Type) Type {
	// this switch is spelled out here, in the Types section:
	// https://github.com/golang/example/tree/master/gotypes
	switch gType.(type) {
	case *types.Basic: return GoBasicToAxePrimative(gType.(*types.Basic))
	case *types.Pointer: return GoPointerToAxeKind(gType.(*types.Pointer))
	case *types.Array: return GoArrayToAxeArray(gType.(*types.Array))
	case *types.Slice: return GoSliceToAxeKind(gType.(*types.Slice))
	case *types.Map: return GoMapToAxeKind(gType.(*types.Map))
	case *types.Chan: return GoChanToAxeKind(gType.(*types.Chan))
	case *types.Struct: return GoStructToAxeStruct(gType.(*types.Struct))
	case *types.Signature: return GoSignatureToAxeFunction(gType.(*types.Signature))
	// case *types.Named: return GoNamedtoAxeName(gType.(*types.Named))
	default: return &Hole{}
	}
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
		log.Println(fragment.url, GoTypeToAxeType(fragment.typ).String())
	}
}
