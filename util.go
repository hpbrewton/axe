package main

import (
	"sort"
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