package main

import (
	"math"
	"sort"
	// "log"
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

func (pq *priorityQueue) reportUpto(upto float64) []int {
	for i, v := range pq.scores {
		if v > upto {
			return pq.indicies[:i]
		}
	}
	return pq.indicies
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
	position int 
	radius float64
	left *VPTree
	right *VPTree
}

// O(nlogn X m), where m is the cost of the metric
func NewVPTree(data []int, metric func(int, int)float64) *VPTree {
	if len(data) == 0 {
		return nil
	}
	if len(data) == 1 {
		return &VPTree{
			position: data[0],
			radius: 0,
			left: nil,
			right: nil,
		}
	}

	pivot := data[0]
	sort.Slice(data[1:], func(i, j int) bool{
		return metric(pivot, data[i]) < metric(pivot, data[j])
	})

	split := len(data)/2
	radius := metric(pivot, data[split-1])
	sides := make(chan bool, 2)
	var left *VPTree 
	var right *VPTree
	go func(){
		left = NewVPTree(data[1:split], metric)
		sides <- true 
	}()
	go func(){
		right = NewVPTree(data[split:], metric)
		sides <- true 
	}()
	for i := 0; i < 2; i++ {
		<-sides 
	}
	return &VPTree{
		position: pivot,
		radius: radius,
		left: left,
		right: right,
	}
}

func (vpt *VPTree) lookupAux(pq *priorityQueue, metric func(v int)float64, cutoff float64) {
	if vpt == nil {
		return
	}

	distance := metric(vpt.position)

	if distance <= cutoff {
		pq.insert(vpt.position, distance)
	}
	if distance <= vpt.radius + cutoff {
		vpt.left.lookupAux(pq, metric, cutoff)
	} 
	if distance >= vpt.radius - cutoff {
		vpt.right.lookupAux(pq, metric, cutoff)
	}
}

// O(logn X m), where m is the cost of the metric
func (vpt *VPTree) Lookup(metric func(v int)float64, k int, cutoff float64) []int {
	pq := mkPQ(k)
	vpt.lookupAux(pq, metric, cutoff)
	return pq.report() //reportUpto(cutoff)
}
