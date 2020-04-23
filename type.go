package main

import (
	"fmt"
	"strings"
	"strconv"
	"reflect"
	"sort"
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

func SortedMethods(methodMap map[string]*Function) []Type {
	
	methodNames := make([]string, len(methodMap))
	methods := make([]interface{}, len(methodMap))
	i := 0
	for name, method := range methodMap {
		methodNames[i] = name 
		methods[i] = interface{}(method)
		i++
	}
	SortBy(sort.StringSlice(methodNames), methods)
	actualMethods := make([]Type, len(methodMap))
	for i, intfMethod := range methods {
		actualMethods[i] = intfMethod.(Type)
	}
	return actualMethods
}

func (mh *MethodHaver) String() string {
	return mh.name
}

type Primative struct {
	Name string
}

func (primative *Primative) String() string {
	return primative.Name
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
	Object Type // hey bud, it's worth noting that this can be nil
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

