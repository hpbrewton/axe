package main

import (
	"testing"
	"reflect"
)

var mc = &MetricConfig{
	AddElemCost: 2.0,
	AddElemCostChan: 10.0,
	DefaultFieldCost: 2.0,
}

func DistanceCheck(t *testing.T, l, r interface{}, expected float64) {
	typ := reflect.TypeOf(l)
	metric, err  := mc.GetMetric(typ)
	if err != nil {
		t.Error(err)
		t.Fatalf("expected to be able to get metric for %v", l)
	}
	d := metric(l, r)
	if expected != d {
		t.Errorf("expected distance %f, but got %f", expected, d)
	}
}

type point struct {
	X int
	Y int
}

type unexpoint struct {
	x int 
	y int
}

func TestDistances(t *testing.T){
	DistanceCheck(t, true, false, 1.0)
	DistanceCheck(t, int(1), int(22), 21.0)
	DistanceCheck(t, int8(1), int8(22), 21.0)
	DistanceCheck(t, int16(1), int16(22), 21.0)
	DistanceCheck(t, int32(1), int32(22), 21.0)
	DistanceCheck(t, int64(1), int64(22), 21.0)
	DistanceCheck(t, uint(1), uint(22), 21.0)
	DistanceCheck(t, uintptr(1), uintptr(22), 1.0)
	DistanceCheck(t, uint8(1), uint8(22), 21.0)
	DistanceCheck(t, uint16(1), uint16(22), 21.0)
	DistanceCheck(t, uint32(1), uint32(22), 21.0)
	DistanceCheck(t, uint64(1), uint64(22), 21.0)
	DistanceCheck(t, float32(1), float32(22), 21.0)
	DistanceCheck(t, float64(1), float64(22), 21.0)
	DistanceCheck(t, complex64(2+5i), complex64(5+9i), 5.0)
	DistanceCheck(t, complex128(2+5i), complex128(5+9i), 5.0)
	DistanceCheck(t, []int{7, 2, 3}, []int{2}, 4.0)
	DistanceCheck(t, point{1, 2}, point{3, 4}, 2.0*((3-1)+(4-2)))
	DistanceCheck(t, unexpoint{1, 2}, unexpoint{3, 4}, 0)
	l := make(chan int, 5)
	l <- 1 
	l <- 2
	close(l)
	r := make(chan int, 5)
	r <- 1
	r <- 14
	close(r)
	DistanceCheck(t, l, r, 12)
	DistanceCheck(t, "cat", "coat", 1.0)
}