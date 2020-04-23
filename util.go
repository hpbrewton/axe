package main

import (
	"sort"
	"math"
)

type SortByer struct {
	l sort.Interface
	r []interface{}
}

func (sb *SortByer) Len() int {
	return sb.l.Len()
}

func (sb *SortByer) Less(i, j int) bool {
	return sb.l.Less(i, j)
}

func (sb *SortByer) Swap(i, j int) {
	sb.l.Swap(i, j)
	sb.r[i], sb.r[j] = sb.r[j], sb.r[i]
} 

func SortBy(l sort.Interface, r []interface{}) {
	sort.Sort(&SortByer{l: l, r : r})
}

func BoolToInt(f bool) int {
	if f {
		return 0
	} else {
		return 1
	}
}

func Min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func Max(a, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}


type Levenshteiner interface {
	Set(i, j int, v float64)
	Get(i, j int) float64 // should return -1 if not set
	Score(i, j int) float64
	Offset() float64
}

func Levenshtein(ler Levenshteiner, i, j int) (v float64) {
	if s := ler.Get(i, j); s >= 0 {
		v = s 
		return
	}

	if Min(i, j) == 0 {
		v = float64(Max(i, j)) * ler.Offset()
	} else {
		idown := Levenshtein(ler, i-1, j) + ler.Offset()
		jdown := Levenshtein(ler, i, j-1) + ler.Offset()
		ddown := Levenshtein(ler, i-1, j-1) + ler.Score(i, j)
		v = math.Min(math.Min(idown, jdown), ddown)
	}
	ler.Set(i, j, v)
	return
}



