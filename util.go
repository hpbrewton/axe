package main

import (
	"sort"
	"log"
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
}

type RuneLevenshteiner struct {
	l []rune
	r []rune 
	store [][]float64
	saved int
}

func NewRuneLevenshteiner(l []rune, r []rune) *RuneLevenshteiner {
	store := make([][]float64, len(l)+1)
	for i, _ := range store {
		store[i] = make([]float64, len(r)+1)
		for j, _ := range store[i] {
			store[i][j] = -1
		}
	}
	return &RuneLevenshteiner {
		l: l, 
		r: r, 
		store: store,
		saved: 0,
	}
}

func (rl *RuneLevenshteiner) Set(i, j int, v float64) {
	rl.store[i][j] = v
}

func (rl *RuneLevenshteiner) Get(i, j int) float64 {
	if rl.store[i][j] >= 0 {
		rl.saved++
	}
	return rl.store[i][j]
}

func (rl *RuneLevenshteiner) Score(i, j int) float64 {
	if i < 0 && j >= 0 {
		return 1 
	} else if i >= 0 && j < 0 {
		return 1
	} else if i < 0 && j < 0 {
		return 0
	} else if rl.l[i-1] == rl.r[j-1] {
		return 0
	} else if rl.l[i-1] == 'i' && rl.r[j-1] == 'j' {
		return 0.5
	} else {
		return 1
	}
}

func Levenshtein(ler Levenshteiner, i, j int) (v float64) {
	log.Println(ler.(*RuneLevenshteiner).store)
	if s := ler.Get(i, j); s >= 0 {
		v = s 
		return
	}

	if Min(i, j) == 0 {
		v = float64(Max(i, j))
	} else {
		idown := Levenshtein(ler, i-1, j) + 1
		jdown := Levenshtein(ler, i, j-1) + 1
		ddown := Levenshtein(ler, i-1, j-1) + ler.Score(i, j)
		v = math.Min(math.Min(idown, jdown), ddown)
	}
	ler.Set(i, j, v)
	return
}

func main() {
	a := "harrison"
	b := "haj"
	rl := NewRuneLevenshteiner([]rune(a), []rune(b))
	d := Levenshtein(rl, len(a), len(b))
	log.Println(d, rl.saved, rl.store)
}