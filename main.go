package main 

import (
	"fmt"
	"sort"
)

func main() {
	fragments, err := GoFragmentsFromDirectory("data/go/src")
	if err != nil {
		panic(err)
	}

	scores := make([]float64, len(fragments))
	query := &Function{
		object: nil,
		arguments: []Type{&MethodHaver{
			name: "Reader",
			self: &Hole{},
			methods: make(map[string]*Function),
		}},
		output: []Type{},
	}
	for i, fragment := range fragments {
		scores[i] = DistanceFragment(fragment, query)
	}
	sort.Sort(&FragmentStoreWithScore{fragments: fragments, scores: scores})

	for i, fragment := range fragments {
		fmt.Println(scores[i], fragment.url, fragment.typ)
	}
}