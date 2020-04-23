package main

import (
	"testing"
	"reflect"
)


var mc = &MetricConfig{
	AddElemCost: 15.0,
	AddElemCostChan: 10.0,
	DefaultFieldCost: 1.0,
	FieldMultipliers: make(map[string]*float64),
	InterfaceCap: 12.0,
	Hole: reflect.TypeOf(&Hole{}),
	HoleDistance: 0.0001,
}

func TestPrimative(t *testing.T) {
	l := &Primative{Name: "int32"}
	r := &Primative{Name: "int64"}
	m, err := mc.GetMetric(reflect.TypeOf(l), reflect.TypeOf(r))
	if err != nil {
		t.Fatalf("%s", err)
	}
	d := m(l, r)
	if d != 2 {
		t.Fatalf("expected 2, but got %f", d)
	}
}


func TestFunction(t *testing.T) {
	l := &Function{
		Object: nil,
		Arguments: []Type{
			&Primative{Name: "int32"},
			&Hole{},
		},
		Output: []Type{
			&Primative{Name: "bool"},
		},
	}
	r := &Function{
		Object: nil,
		Arguments: []Type{
			&Primative{Name: "int64"},
			&Primative{Name: "int32"},
		},
		Output: []Type{
			&Primative{Name: "bool"},
		},
	}
	m, err := mc.GetMetric(reflect.TypeOf(l), reflect.TypeOf(r))
	if err != nil {
		t.Fatalf("%s", err)
	}
	if m(l, r) != 2+mc.HoleDistance {
		t.Fatalf("should have been ze same")
	}
}