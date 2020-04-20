package main

import (
	"os"
	"testing"
	"fmt"
	"bufio"
)

const (
	googleEnglish10000Path = "data/google-10000-english-usa.txt"
)

func GoogleEnglish10000() [][]rune {
	file, err := os.Open(googleEnglish10000Path)
	if err != nil {
		panic(fmt.Sprintf("need file %s", googleEnglish10000Path))
	}
	breader := bufio.NewReader(file)

	lines := make([][]rune, 10000)
	for i := 0; i < 10000; i++ {
		str, err := breader.ReadString('\n')
		lines[i] = []rune(str)
		lines[i] = lines[i][:len(lines[i])-1]
		if err != nil {
			panic(fmt.Sprintf("problem in reading %s", googleEnglish10000Path))
		}
	}

	return lines
}

func equal(l, r []int) bool {
	if len(l) != len(r) {
		return false 
	} else {
		for i, v := range l {
			if v != r[i] {
				return false 
			}
		}
		return true
	}
}

func equalRuneSlice(l, r []rune) bool {
	if len(l) != len(r) {
		return false 
	} else {
		for i, v := range l {
			if v != r[i] {
				return false 
			}
		}
		return true
	}
}

func TestPriority(t *testing.T) {
	pq := mkPQ(3)
	for i, v := range []float64{12, 4, 5, 13, 2, 54} {
		pq.insert(i, v)
	} 
	expected := []int{4, 1, 2}
	if !equal(expected, pq.report()) {
		t.Errorf("expected final indicies %v but got %v", expected, pq.indicies)
	}

	pq = mkPQ(2)
	for i, v := range []float64{1, 2} {
		pq.insert(i, v)
	}
	expected = []int{0, 1}
	if !equal(expected, pq.report()) {
		t.Errorf("expected final indicies %v but got %v", expected, pq.indicies)
	}
}

func Looker(objects [][]rune, key []rune) func(v int)float64 {
	return func(v int)float64{
		other := objects[v]
		rl := NewRuneLevenshteiner(other, key)
		d := Levenshtein(rl, len(other), len(key))
		return d
	}
}

func TestVPTree(t *testing.T) {
	vpt := NewVPTree()
	words := GoogleEnglish10000()
	for i, _ := range words{
		vpt.Insert(i)
	}
	metric := Looker(words, []rune("cat"))
	
	output := []string{
		"cat",
		"at",
		"can",
		"car",
		"cart",
	}

	for i, out := range vpt.Lookup(metric, 5) {
		strout := string(words[out])
		if strout != output[i]{
			t.Errorf("expected %s, but got %s", output[i], strout)
		}
	}
}