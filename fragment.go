package main

import "go/types"

type GoFragment struct {
	comment string
	pkg *types.Package
	url  string
	typ  Type
	line int
	col  int
}

func DistanceFragment(l *GoFragment, q Query) float64 {
	dc := DefaultDistanceConfig
	return dc.Distance(l.typ, q.(Type))
}

type FragmentStore []*GoFragment 

func (fs FragmentStore) Len() int {
	return len(fs ) 
}

func (fs FragmentStore) Less(i, j int) bool {
	return Ord(fs[i].typ, fs[j].typ) > 0 
}
 
func (fs FragmentStore) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

type FragmentStoreWithScore struct {
	fragments []*GoFragment 
	scores []float64
}

func (fsws *FragmentStoreWithScore) Len() int {
	return len(fsws.scores)
}

func (fsws *FragmentStoreWithScore) Less(i, j int ) bool {
	return fsws.scores[i] < fsws.scores[j]
}

func (fsws *FragmentStoreWithScore) Swap(i, j int) {
	fsws.scores[i], fsws.scores[j] = fsws.scores[j], fsws.scores[i]
	fsws.fragments[i], fsws.fragments[j] = fsws.fragments[j], fsws.fragments[i]
}
