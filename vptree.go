package main

import (
	"math"
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

func (pq *priorityQueue) report() []int {
	return pq.indicies[:pq.seen]
}

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
	Store []int // a list of indicies
}

func NewVPTree() *VPTree {
	return &VPTree {
		Store: make([]int, 1000),
	}
}

func (vpt *VPTree) Insert(v int) {
	vpt.Store = append(vpt.Store, v)
}

func (vpt *VPTree) Lookup(metric func(v int)float64, k int) []int {
	pq := mkPQ(k)
	for _, v := range vpt.Store {
		pq.insert(v, metric(v))
	}
	return pq.report()
}