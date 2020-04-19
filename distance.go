package main

import (
	"math"
	"reflect"
	"strings"
)

const (
	Same = 0
	Unknown = 1000
)

var Distant = math.Inf(1)

var DefaultDistanceConfig = &DistanceConfig {
	changeOrder: 0,
	arguments: 2,
	returns: 1,
	recieverChange: 0.5,
	recieverAdd: 0.75,
	structField: 2,
	abstactUp: 0.1,
	primativeChange: 0.1,
	tupleFromNothing: 10,
	arraySizeSpecification: 5,
	arraySizeChange: 10,
	differentKind: 10,
	kindArgsDiffernt: 1,
	differentObjectName: 0.1,
	differentMethods: 0.9,
	differtObject: 0.2,
	differentInterfaceName: 0.2,
	differentImplements: 5,
	differentMethodsInterface: 5,
	addVar: 0.2,
}

type DistanceConfig struct {
	changeOrder float64
	arguments float64
	returns float64
	recieverChange float64
	recieverAdd float64
	structField float64
	abstactUp float64 // go to interface from an implementer
	primativeChange float64
	tupleFromNothing float64
	arraySizeSpecification float64 // []obj <-> [n]obj
	arraySizeChange float64 // [n]obj <-> [m] obj
	arrayTypeChange float64 // [?]obj <-> [?]jbo
	differentKind float64 // chan obj <-> [] obj // this might be better with specific kind changes by language but that is a bit out of scope for now :(
	kindArgsDiffernt float64 // []obj <-> []jbo // 
	differentObjectName float64
	differentMethods float64
	differtObject float64
	differentInterfaceName float64
	differentMethodsInterface float64
	differentImplements float64
	addVar float64
}

func (dc *DistanceConfig) DistancePrimative(l, r *Primative) float64 {
	if strings.Compare(l.name, r.name) == 0 {
		return Same
	} else {
		return dc.primativeChange
	}
}

func (dc *DistanceConfig) DistanceTuple(l, r []Type) (score float64) {
	s := make([][]float64, len(l)+1)
	for i, _ := range s {
		s[i] = make([]float64, len(r)+1)
		for j, _ := range s[i] {
			s[i][j] = -1
		}
	}
	tl := &TypeLevenshteiner {
		l: l,
		r: r,
		s: s,
		dc: dc,
	}
	score = Levenshtein(tl, len(l), len(r))
	return
}

func (dc *DistanceConfig) DistanceInterface(l, r *Interface) (score float64) {
	if l.name != r.name {
		score += dc.differentInterfaceName
	}

	leftMethods := SortedMethods(l.methods)
	rightMethods := SortedMethods(r.methods)
	score += dc.differentMethodsInterface*dc.DistanceTuple(leftMethods, rightMethods)

	score += dc.differentImplements*dc.DistanceTuple(l.implements, r.implements)
	return 
}

func (dc *DistanceConfig) DistanceMethodHaver(l, r *MethodHaver) (score float64) {
	if l.name != r.name {
		score += dc.differentObjectName 
	}

	score += dc.differtObject*dc.Distance(l.self, r.self) 
	leftMethods := SortedMethods(l.methods)
	rightMethods := SortedMethods(r.methods)
	score += dc.differentMethods*dc.DistanceTuple(leftMethods, rightMethods)

	return
}

func (dc *DistanceConfig) DistanceKind(l, r *Kind) (score float64) {
	if l.name != r.name {
		score += dc.differentKind
	}

	score += dc.kindArgsDiffernt * dc.DistanceTuple(l.arguments, r.arguments)

	return 
}

func (dc *DistanceConfig) DistanceArray(l, r *Array) (score float64) {
	if l.size < 0 && r.size >= 0 || l.size >= 0 && r.size < 0 {
		score += dc.arraySizeSpecification
	} else if l.size != r.size {
		score += dc.arraySizeChange
	}

	score += dc.arrayTypeChange * dc.Distance(l.typ, r.typ)
	return 
}

func (dc *DistanceConfig) DistanceFunction(l, r *Function) (score float64) {
	// first we compare the recievers, see if they are the same 
	if l.object != r.object {
		if l.object == nil || r.object == nil {
			score += dc.recieverAdd 
		} else {
			score += dc.recieverChange * dc.Distance(l.object, r.object)
		}
	} 

	score += dc.arguments * dc.DistanceTuple(l.arguments, r.arguments)
	score += dc.returns * dc.DistanceTuple(l.output, r.output)

	return 
}

func (dc *DistanceConfig) DistanceStruct(l, r *Struct) (score float64) {
	score += dc.structField * dc.DistanceTuple(l.fields, r.fields)

	return 
}

func (dc *DistanceConfig) Distance(l Type, r Type) float64 {
	if IsHole(l) || IsHole(r) {
		return Same 
	}

	if reflect.DeepEqual(l, r) {
		return Same
	}

	if reflect.TypeOf(l) != reflect.TypeOf(r) { // this is not right, but an approx
		return Distant
	}

	switch l.(type) {
	case *Primative: return dc.DistancePrimative(l.(*Primative), r.(*Primative))
	case *Interface: return dc.DistanceInterface(l.(*Interface), r.(*Interface))
	case *MethodHaver: return dc.DistanceMethodHaver(l.(*MethodHaver), r.(*MethodHaver))
	case *Kind: return dc.DistanceKind(l.(*Kind), r.(*Kind))
	case *Array: return dc.DistanceArray(l.(*Array), r.(*Array))
	case *Function: return dc.DistanceFunction(l.(*Function), r.(*Function))
	case *Struct: return dc.DistanceStruct(l.(*Struct), r.(*Struct))
	default:
		return Unknown
	}
}