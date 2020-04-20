package main

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