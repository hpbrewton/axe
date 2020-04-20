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

func TestIndex(t *testing.T) {
	data := []int{64, 32, 16, 8, 4, 15}
	indicies := make([]int, len(data))
	for i, _ := range data {
		indicies[i] = i
	}
	radii := index(indicies, func(a, b int) float64{
		together := data[a] ^ data[b]
		count := 0.0
		for together != 0 {
			if together%2 != 0 {
				count += 1.0
			}
			together/=2
		}
		return count
	})
	t.Log(radii)
}

func TestVPTree(t *testing.T) {
	words := GoogleEnglish10000()
	indicies := make([]int, 10000)
	for i, _ := range words{
		indicies[i] = i
	}
	vpt := NewVPTree(indicies, func(a int, b int) float64 {
		aword := words[a]
		bword := words[b]
		rl := NewRuneLevenshteiner(aword, bword)
		return Levenshtein(rl, len(aword), len(bword))
	})

	metric := Looker(words, []rune("cat"))
	output := []string{
		"cat",
		"mat",
		"rat",
		"chat",
		"sat",
	}
	for i, out := range vpt.Lookup(metric, 5, 1.0) {
		strout := string(words[out])
		if strout != output[i]{
			t.Errorf("expected %s, but got %s", output[i], strout)
		}
	}
}