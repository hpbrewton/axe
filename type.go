package main

import (
	"fmt"
	"strings"
	"strconv"
)

type Type interface {
	fmt.Stringer
}

type Named struct {
	name   string
	actual Type
}

func (named *Named) String() string {
	return named.name
}

type Primative struct {
	name string
}

func (primative *Primative) String() string {
	return primative.name
}

type Kind struct {
	name string 
	arguments []Type
}

func (kind *Kind) String() string {
	return fmt.Sprintf("(%s %v)", kind.name, kind.arguments)
}

// a negative indicates unknown size
type Array struct {
	size int
	typ Type
}

func (array *Array) String() string {
	var sizeStr string
	if array.size < 0 {
		sizeStr = ""
	} else {			
		sizeStr = strconv.Itoa(array.size)
	}
	return fmt.Sprintf("[%s]%s", sizeStr, array.typ)
}

type Function struct {
	arguments []Type
	output []Type
}

func (function *Function) String() string {
	return fmt.Sprintf("%v -> %v", function.arguments, function.output)
}

type Struct struct {
	fieldNames []string
	fields []Type 
}

func (strct *Struct) String() string {
	var b strings.Builder
	b.WriteString("<")
	for i, fieldName := range strct.fieldNames {
		fmt.Fprintf(&b, "%s:%v, ", fieldName, strct.fields[i])
	}
	b.WriteString(">")
	return b.String()
}

type Hole struct{}

func (*Hole) String() string {
	return "??"
}

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
