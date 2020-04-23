package main 

import (
	"fmt"
	"reflect"
	// "runtime/debug"
)


var mc = &MetricConfig{
	AddElemCost: 1.0,
	AddElemCostChan: 10.0,
	DefaultFieldCost: 1.0,
	FieldMultipliers: make(map[string]*float64),
	InterfaceCap: 1000.0,
	Hole: reflect.TypeOf(&Hole{}),
	HoleDistance: 0.0001,
}

func main() {
	fragments, err := GoFragmentsFromDirectory("data/go/src")
	if err != nil {
		panic(err)
	}
	indicies := make([]int, len(fragments))
	for i, _ := range fragments {
		indicies[i] = i
	}
	fragType := reflect.TypeOf(fragments[0])
	fragMetric, err := mc.GetMetric(fragType, fragType)
	if err != nil {
		fmt.Errorf("%s", err)
	}
	metric := func(a, b int) float64 {
		// fmt.Println(a, b, fragments[a].Typ, fragments[b].Typ)
		return fragMetric(fragments[a], fragments[b])
	}
	vpt := NewVPTree(indicies, metric)
	query := &GoFragment{
		Comment: "hooplah",
		pkg: nil,
		url: "",
		Typ: &Function{
			Object: nil,
			Arguments: []Type{&Kind{
				Name: "*",
				Arguments: []Type{&Primative{
					Name: "string",
			}}}},
			Output: []Type{
				&Hole{},
				&Hole{}, 
			},
		},
	}
	looker := func(a int) float64{
		return fragMetric(query, fragments[a])
	}
	recieved := vpt.Lookup(looker, 10, 1000.0)
	for _, recv := range recieved {
		fmt.Println(fragments[recv].Typ)
	}
}