package main 

type TypeLevenshteiner struct {
	l []Type
	r []Type
	s [][]float64
	dc *DistanceConfig
}

func (tl *TypeLevenshteiner) Set(i, j int, v float64) {
	tl.s[i][j] = v 
} 

func (tl *TypeLevenshteiner) Get(i, j int) float64 {
	return tl.s[i][j]
}

func (tl *TypeLevenshteiner) Score(i, j int) float64 {
	if i < 0 && j >= 0 {
		return tl.dc.addVar
	} else if i >= 0 && j < 0 {
		return tl.dc.addVar
	} else if i < 0 && j < 0 {
		return Same
	} else {
		return tl.dc.Distance(tl.l[i-1], tl.r[j-1])
	}
}