package main

import (
	"fmt"
	"math"
	"reflect"
	"errors"
	"log"
	"math/cmplx"
)

func toArrayEmptyInterface(ev interface{})[]interface{} {
	v := reflect.ValueOf(ev)
	length := v.Len()
	ret := make([]interface{}, length)
	for i, _ := range ret {
		ret[i] = v.Index(i).Interface()
	}
	return ret
}

func toArrayEmptyInterfaceChan(ev interface{})[]interface{} {
	v := reflect.ValueOf(ev)
	ret := make([]interface{}, 0)
	for {
		if x, ok := v.Recv(); ok {
			ret = append(ret, x.Interface())
		} else {
			break
		}
	}
	return ret
}

type MetricConfig struct {
	AddElemCost float64
	AddElemCostChan float64
	DefaultFieldCost float64
}

type Metric func(a, b interface{}) float64 

func (mc *MetricConfig) StringMetric() Metric {
	return func(l, r interface{}) float64 {
		leftRunes := []rune(l.(string))
		rightRunes := []rune(r.(string))
		rl := NewRuneLevenshteiner(leftRunes, rightRunes)
		score := Levenshtein(rl, len(leftRunes), len(rightRunes))
		return score
	}
}

type Numberer func(interface{})float64

func fromBool(v interface{})float64{
	if v.(bool) {
		return 1
	} else {
		return 0
	}
}
func fromInt(v interface{})float64{ return float64(v.(int))}
func fromInt8(v interface{})float64{ return float64(v.(int8))}
func fromInt16(v interface{})float64{ return float64(v.(int16))}
func fromInt32(v interface{})float64{ return float64(v.(int32))}
func fromInt64(v interface{})float64{ return float64(v.(int64))}
func fromUint(v interface{})float64{ return float64(v.(uint))}
func fromUint8(v interface{})float64{ return float64(v.(uint8))}
func fromUint16(v interface{})float64{ return float64(v.(uint16))}
func fromUint32(v interface{})float64{ return float64(v.(uint32))}
func fromUint64(v interface{})float64{ return float64(v.(uint64))}
func fromFloat32(v interface{})float64{ return float64(v.(float32))}
func fromFloat64(v interface{})float64{ return float64(v.(float64))}
func (mc *MetricConfig) NumberMetric(converter Numberer) Metric {
	return func(l, r interface{}) float64 {
		nl, nr := converter(l), converter(r)
		return math.Abs(nl - nr)
	}
}
func (mc *MetricConfig) DiscreteMetric(l, r interface{}) float64 {
	if l == r {
		return 0.0
	} else {
		return 1.0
	}
}
func (mc *MetricConfig) Complex64Metric(l, r interface{}) float64 {
	return cmplx.Abs(complex128(l.(complex64)) - complex128(r.(complex64)))
}
func (mc *MetricConfig) Complex128Metric(l, r interface{}) float64 {
	return cmplx.Abs(l.(complex128) - r.(complex128))
}
func (mc *MetricConfig) ArrayMetric(e reflect.Type) (Metric, error) {
	typ := e.Elem()
	elemMetric, err := mc.GetMetric(typ)
	if err != nil {
		return nil, err 
	} else {
		f := func(l, r interface{})float64{
			la := toArrayEmptyInterface(l)
			ra := toArrayEmptyInterface(r)
			sl := NewSliceLevenshteiner(mc.AddElemCost, elemMetric, la, ra)
			return Levenshtein(sl, len(la), len(ra))
		}
		return f, nil
	}
}
func (mc *MetricConfig) ChanMetric(e reflect.Type) (Metric, error) {
	typ := e.Elem()
	elemMetric, err := mc.GetMetric(typ)
	if err != nil {
		return nil, err 
	} else {
		f := func(l, r interface{})float64{
			la := toArrayEmptyInterfaceChan(l)
			ra := toArrayEmptyInterfaceChan(r) 
			sl := NewSliceLevenshteiner(mc.AddElemCostChan, elemMetric, la, ra)
			d := Levenshtein(sl, len(la), len(ra))
			log.Println(sl.store)
			return d
		}
		return f, nil
	}
}
func (mc *MetricConfig) StructMetric(t reflect.Type) (Metric, error) {
	nfields := t.NumField()

	f := func (l, r interface{})float64 {
		var total float64 = 0 
		lo, ro := reflect.ValueOf(l), reflect.ValueOf(r)
		for i := 0; i < nfields; i++ {
			structField := t.Field(i)
			if structField.PkgPath != "" {
				continue //not exported, ignore
			}
			fieldType := structField.Type 
			metric, err := mc.GetMetric(fieldType)
			if err != nil {
				panic(err) // we eventually will get rid of all errors
			}
			total += mc.DefaultFieldCost*metric(lo.Field(i).Interface(), ro.Field(i).Interface())
		}
		return total
	}
	return f, nil
}
 
func (mc *MetricConfig) GetMetric(t reflect.Type) (Metric, error) {
	switch t.Kind() {
	case reflect.Bool: return mc.NumberMetric(fromBool), nil
	case reflect.Int: return mc.NumberMetric(fromInt), nil 
	case reflect.Int8: return mc.NumberMetric(fromInt8), nil
	case reflect.Int16: return mc.NumberMetric(fromInt16), nil
	case reflect.Int32: return mc.NumberMetric(fromInt32), nil
	case reflect.Int64: return mc.NumberMetric(fromInt64), nil
	case reflect.Uint: return mc.NumberMetric(fromUint), nil
	case reflect.Uint8: return mc.NumberMetric(fromUint8), nil
	case reflect.Uint16: return mc.NumberMetric(fromUint16), nil
	case reflect.Uint32: return mc.NumberMetric(fromUint32), nil
	case reflect.Uint64: return mc.NumberMetric(fromUint64), nil
	case reflect.Uintptr: return mc.DiscreteMetric, nil
	case reflect.Float32: return mc.NumberMetric(fromFloat32), nil
	case reflect.Float64: return mc.NumberMetric(fromFloat64), nil
	case reflect.Complex64: return mc.Complex64Metric, nil
	case reflect.Complex128: return mc.Complex128Metric, nil
	case reflect.Array: return mc.ArrayMetric(t)
	case reflect.Chan: return mc.ChanMetric(t)
	case reflect.Struct: return mc.StructMetric(t)
	case reflect.String: return mc.StringMetric(), nil 
	case reflect.Slice: return mc.ArrayMetric(t)
	default: return nil, errors.New(fmt.Sprintf("undefined for type %s", t.String()))
	}
}