package main

import "go/types"

type GoFragment struct {
	Comment string `type2vec:"0"`
	pkg *types.Package 
	url  string 
	Typ  Type `type2vec:"1"`
	line int 
	col  int 
}