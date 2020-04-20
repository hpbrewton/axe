package main

import (
	"testing"
)

func TestIJSim(t *testing.T) {
	a := "harrison"
	b := "haj"
	rl := NewRuneLevenshteiner([]rune(a), []rune(b))
	d := Levenshtein(rl, len(a), len(b))
	if d != 5.5 {
		t.Errorf("%s and %s should have adjusted levenshtein distance %f but got %f", a, b, 5.5, d)
	}
}

func TestEmptyIsZero(t *testing.T) {
	a := ""
	b := "haj"
	rl := NewRuneLevenshteiner([]rune(a), []rune(b))
	d := Levenshtein(rl, len(a), len(b))
	if d != 3 {
		t.Errorf("%s and %s should have adjusted levenshtein distance %f but got %f", a, b, 3.0, d)
	}
}