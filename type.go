package main

import (
	"fmt"
)

type Type interface {
}

type Named struct {
	name string
	actual Type
}

type Primative struct {
	name string
}

type Hole struct {}

func Debug(t Type) string {
	switch t.(type) {
	case *Named:
		named := t.(*Named)
		return fmt.Sprintf("%s := %s", named.name, Debug(named.actual))
	case *Primative:
		primative := t.(*Primative)
		return fmt.Sprintf("%s", primative.name)
	default:
		panic(fmt.Sprintf("%T is not a supported type", t))
	}
}
