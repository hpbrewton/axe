package main

import "go/types"

type Fragment struct {
	comment string
	url  string
	typ  types.Type
	line int
	col  int
}
