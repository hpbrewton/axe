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
	name string
	implements []Type
	methods map[string]*Function
}

func (ntrf *Interface) String() string {
	return fmt.Sprintf("%s", ntrf.name)
}

type MethodHaver struct {
	name string
	self Type
	methods map[string]*Function
}

func (mh *MethodHaver) String() string {
	return mh.name
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
	object Type // hey bud, it's worth noting that this can be nil
	arguments []Type
	output []Type
}

func (function *Function) String() string {
	var str string
	if function.object == nil {
		str = ""
	} else {
		str = function.object.String()
		str = fmt.Sprintf("%s -> ", str)
	}
	return fmt.Sprintf("%s%v -> %v", str, function.arguments, function.output)
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

// it is useful to have a canonical ordering of types
// so here it is!
func Ord(a Type,  b Type) int {
	atyp := reflect.TypeOf(a)
	btyp := reflect.TypeOf(b)

	if atyp.Name() == btyp.Name() {
		return 0
	} else if atyp.Name() < btyp.Name() {
		return 1 
	} else {
		return -1
	}
}

