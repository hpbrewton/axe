package main

import (
	"fmt"
	"strings"
	"strconv"
	"reflect"
)

type Type interface {
	fmt.Stringer
}

type Interface struct {
	Name string
	Implements []Type
	methods map[string]*Function
}

func (ntrf *Interface) String() string {
	return fmt.Sprintf("%s", ntrf.Name)
}

type MethodHaver struct {
	Name string
	self Type
	methods map[string]*Function
}

func (mh *MethodHaver) String() string {
	return mh.Name
}

type Primative struct {
	Name string
}

func (primative *Primative) String() string {
	return primative.Name
}

type Kind struct {
	Name string 
	Arguments []Type
}

func (kind *Kind) String() string {
	return fmt.Sprintf("(%s %v)", kind.Name, kind.Arguments)
}

// a negative indicates unknown size
type Array struct {
	Size int
	Typ Type
}

func (array *Array) String() string {
	var sizeStr string
	if array.Size < 0 {
		sizeStr = ""
	} else {			
		sizeStr = strconv.Itoa(array.Size)
	}
	return fmt.Sprintf("[%s]%s", sizeStr, array.Typ)
}

type Function struct {
	Object Type 
	Arguments []Type
	Output []Type
}

func (function *Function) String() string {
	var str string
	if function.Object == nil {
		str = ""
	} else {
		str = function.Object.String()
		str = fmt.Sprintf("%s -> ", str)
	}
	return fmt.Sprintf("%s%v -> %v", str, function.Arguments, function.Output)
}

type Struct struct {
	FieldNames []string
	Fields []Type 
}

func (strct *Struct) String() string {
	var b strings.Builder
	b.WriteString("<")
	for i, fieldName := range strct.FieldNames {
		fmt.Fprintf(&b, "%s:%v, ", fieldName, strct.Fields[i])
	}
	b.WriteString(">")
	return b.String()
}

type Hole struct{}

func (*Hole) String() string {
	return "??"
}

func IsHole(t Type) bool {
	switch t.(type) {
	case *Hole: return true
	default: return false
	}
}

// it is useful to have a canonical ordering of types
// so here it is!
func Ord(a Type,  b Type) int {
	atyp := reflect.TypeOf(a)
	btyp := reflect.TypeOf(b)

	return strings.Compare(atyp.String(), btyp.String())
}

