package main 

type SliceLevenshteiner struct {
	l []interface{}
	r []interface{}
	store [][]float64
	addVarCost float64
	metric Metric
}

func NewSliceLevenshteiner(addVarCost float64, metric Metric, l, r []interface{}) *SliceLevenshteiner {
	store := make([][]float64, len(l)+1)
	for i, _ := range store {
		store[i] = make([]float64, len(r)+1)
		for j, _ := range store[i] {
			store[i][j] = -1
		}
	}
	return &SliceLevenshteiner{
		l: l,
		r: r,
		store: store,
		addVarCost: addVarCost,
		metric: metric,
	}
}

func (sl *SliceLevenshteiner) Set(i, j int, v float64) {
	sl.store[i][j] = v
}

func (sl *SliceLevenshteiner) Get(i, j int) float64 {
	return sl.store[i][j]
}

func (sl *SliceLevenshteiner) Score(i, j int) float64 {
	return sl.metric(sl.l[i-1], sl.r[j-1])
}

func (sl *SliceLevenshteiner) Offset() float64 {
	return sl.addVarCost
}