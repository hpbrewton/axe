package main 

import (
	"fmt"
	"reflect"
	"math/rand"
	"time"
	"os"
	"path/filepath"
	"sync"
	"encoding/json"
	// "runtime/debug"
)


var mc = &MetricConfig{
	AddElemCost: 1.0,
	AddElemCostChan: 10.0,
	DefaultFieldCost: 1.0,
	FieldMultipliers: &sync.Map{},
	InterfaceCap: 1000.0,
	Hole: reflect.TypeOf(&Hole{}),
	HoleDistance: 0.0001,
}

var rando *rand.Rand
var nrepeat int 

func init() {
	rando = rand.New(rand.NewSource(10))
	nrepeat = 10
}

func fragments(dir string) []*GoFragment {
	fragments, err := GoFragmentsFromDirectory(dir)
	// funcFragmentsmake([]*GoFragment, 0)
	// for _, frag := range fragments {
	// 	switch frag.Typ.(type) {
	// 	case *Function:
	// 		fmt.Println(frag.Typ.(*Function))
	// 	}
	// }
	if err != nil {
		panic(err)
	}
	return fragments
}

func mkVPTree(fragMetric Metric, fragments []*GoFragment) (*VPTree, int64) {
	start := time.Now()
	if len(fragments) == 0 {
		return nil, 0
	}
	indicies := make([]int, len(fragments))
	for i, _ := range fragments {
		indicies[i] = i
	}
	metric := func(a, b int) float64 {
		return fragMetric(fragments[a], fragments[b])
	}
	vpt := NewVPTree(indicies, metric)
	duration := time.Now().Sub(start).Nanoseconds()
	return vpt, duration
}

func directoriesOfInterest() []string {
	directoriesOfInterest := make([]string, 0)
	filepath.Walk("data/go/src", func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsDir(){
			gofiles, err := filepath.Glob(fmt.Sprintf("%s/*.go", path))
			if err != nil {
				panic(err)
			} 
			if len(gofiles) > 0 {
				directoriesOfInterest = append(directoriesOfInterest, path)
			}
		}
		return nil
	})
	return directoriesOfInterest
}

type TestResult struct {
	Dir string
	Nfrags int 
	IndexNS int64 
	DataNS []int64
}

func main() {
	fragType := reflect.TypeOf(&GoFragment{})
	fragMetric, err := mc.GetMetric(fragType, fragType)
	if err != nil {
		panic(err)
	}
	for _, directory := range directoriesOfInterest() {
		frags := fragments(directory)
		if len(frags) == 0 {
			continue
		}
		vpt, indexns := mkVPTree(fragMetric, frags)
		data := make([]int64, 0)
		queries := make([]*GoFragment, nrepeat)
		for i, _ := range queries {
			queries[i] = frags[rando.Intn(len(frags))]
		}
		for _, cutoff := range []float64{0, 1, 2, 3, 5, 10, 100, 1000} {
			if err != nil {
				continue
			}
			avg := int64(0)
			for i := 0; i < nrepeat; i++ {
				query := queries[i]
				looker := func(a int) float64{
					return fragMetric(query, frags[a])
				}
				sum := int64(0)
				for j := 0; j < 10; j++ {
					start := time.Now()
					vpt.Lookup(looker, j, cutoff)
					sum += time.Now().Sub(start).Nanoseconds()
				}
				sum /= 10
				avg += sum
			}
			avg /= int64(nrepeat)
			data = append(data, avg)
		}
		testResult := TestResult{
			Dir: directory,
			Nfrags: len(frags),
			IndexNS: indexns,
			DataNS: data,
		}
		jsonResult, err := json.Marshal(testResult)
		if err != nil {
			panic(err)
		}
		os.Stdout.Write(jsonResult)
	}
}