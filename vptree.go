package main

import (
	"math"
	"sort"
)

type priorityQueue struct {
	seen int
	size int
	scores []float64
	indicies []int 
}

func mkPQ(size int) *priorityQueue {
	scores := make([]float64, size)
	indicies := make([]int, size)
	for i, _ := range scores {
		scores[i] = math.Inf(1)
	}
	return &priorityQueue{
		seen: 0,
		size: size, 
		scores: scores,
		indicies: indicies,
	}
}

// O(1)
func (pq *priorityQueue) report() []int {
	return pq.indicies[:pq.seen]
}

// O(|pq|)
func (pq *priorityQueue) insert(k int, s float64) {
	pq.indicies = append(pq.indicies, k)
	pq.scores = append(pq.scores, s)
	for i := pq.size; i > 0; i-- {
		if pq.scores[i] < pq.scores[i-1] {
			pq.scores[i], pq.scores[i-1] = pq.scores[i-1], pq.scores[i]
			pq.indicies[i], pq.indicies[i-1] = pq.indicies[i-1], pq.indicies[i]
		} else {
			goto end 
		}
	}
end:
	if pq.seen < pq.size {
		pq.seen++
	}
	pq.indicies = pq.indicies[:pq.size]
	pq.scores = pq.scores[:pq.size]
}

type VPTree struct {
	radii []float64
	positions []int // a list of indicies
}

// O(nlogn X m), where m is the cost of the metric
func NewVPTree(data []int, metric func(int, int)float64) *VPTree {
	radii := index(data, metric)
	return &VPTree {
		radii: radii,
		positions: data,
	}
}

// O(nlogn X m), where m is the cost of the metric
func index(data []int, metric func(int, int)float64) []float64{
	if len(data) <= 1 {
		return make([]float64, len(data))
	}

	pivot := data[0] // maybe do something fancier here TODO
	sort.Slice(data, func(i, j int) bool{
		return metric(pivot, i) < metric(pivot, j)
	})
	split := 1+(len(data)-1)/2
	radius := metric(pivot, data[split-1]) // the furthest in the circle

	lradii := index(data[1:split], metric)
	rraddi := index(data[split:], metric)
	return append([]float64{radius}, append(lradii, rraddi...)...)
}

func lookupAux(positions []int,
	radii []float64,
	pq *priorityQueue,
	metric func(v int)float64,
	cutoff float64,
	){
	if len(positions) == 0 {
		return
	} else if len(positions) == 1 {
		v := positions[0]
		pq.insert(v, metric(v))
	} else {
		node := positions[0]
		// radius := radii[0]
		distance := metric(node)

		if distance < cutoff {
			pq.insert(node, distance)
		}

		split := 1+(len(positions)-1)/2
		if true { // distance <  radius + cutoff {
			lookupAux(positions[1:split], radii[1:split], pq, metric, cutoff)
		}
		if true { // distance >= radius - cutoff {
			lookupAux(positions[split:], radii[split:], pq, metric, cutoff)
		}
	}
}

// O(logn X m), where m is the cost of the metric
func (vpt *VPTree) Lookup(metric func(v int)float64, k int, cutoff float64) []int {
	pq := mkPQ(k)
	lookupAux(vpt.positions, vpt.radii, pq, metric, cutoff)
	return pq.report()
}
