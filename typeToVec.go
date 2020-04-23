package main

import (
	"fmt"
	"math"
	"reflect"
	"errors"
	"strconv"
	"math/cmplx"
	"log"
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
	FieldMultipliers map[string]*float64
	InterfaceCap float64
	Hole reflect.Type
	HoleDistance float64
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
func (mc *MetricConfig) DiscreteMetric(dist float64) Metric {
	return func(l, r interface{})float64{
		if l == r {
			return 0.0
		} else {
			return dist
		}
	}	
}
func (mc *MetricConfig) Complex64Metric(l, r interface{}) float64 {
	return cmplx.Abs(complex128(l.(complex64)) - complex128(r.(complex64)))
}
func (mc *MetricConfig) Complex128Metric(l, r interface{}) float64 {
	return cmplx.Abs(l.(complex128) - r.(complex128))
}
func (mc *MetricConfig) ArrayMetric(lt, rt reflect.Type) (Metric, error) {
	lelem := lt.Elem()
	relem := rt.Elem()
	elemMetric, err := mc.GetMetric(lelem, relem)
	if err != nil {
		return nil, err 
	} else {
		f := func(l, r interface{})float64{
			la := toArrayEmptyInterface(l)
			ra := toArrayEmptyInterface(r)
			sl := NewSliceLevenshteiner(mc.AddElemCost, elemMetric, la, ra)
			d := Levenshtein(sl, len(la), len(ra))
			log.Println(d)
			return d
		}
		return f, nil
	}
}
func (mc *MetricConfig) ChanMetric(lt, rt reflect.Type) (Metric, error) {
	lelem := lt.Elem()
	relem := rt.Elem()
	elemMetric, err := mc.GetMetric(lelem, relem)
	if err != nil {
		return nil, err 
	} else {
		f := func(l, r interface{})float64{
			la := toArrayEmptyInterfaceChan(l)
			ra := toArrayEmptyInterfaceChan(r) 
			sl := NewSliceLevenshteiner(mc.AddElemCostChan, elemMetric, la, ra)
			d := Levenshtein(sl, len(la), len(ra))
			return d
		}
		return f, nil
	}
}
func (mc *MetricConfig) StructMetric(lt, rt reflect.Type) (Metric, error) {
	nfields := lt.NumField()

	f := func (l, r interface{})float64 {
		var total float64 = 0 
		lo, ro := reflect.ValueOf(l), reflect.ValueOf(r)
		for i := 0; i < nfields; i++ {
			structField := lt.Field(i)
			if structField.PkgPath != "" {
				continue //not exported, ignore
			}
			fieldType := structField.Type 
			tag := structField.Tag

			metric, err := mc.GetMetric(fieldType, rt.Field(i).Type)
			if err != nil {
				panic(err) // we eventually will get rid of all errors
			}

			// set field identifiers if not present, otw get it 
			identifier := fmt.Sprintf("%s/%s", structField.Name, lt.PkgPath())
			var multipler float64
			if val, ok := mc.FieldMultipliers[identifier]; ok {
				multipler = *val 
			} else {
				field := tag.Get("type2vec")
				foundMultiplier, err := strconv.ParseFloat(field, 64)
				if err != nil {
					multipler = mc.DefaultFieldCost
				} else {
					multipler = foundMultiplier
				}
				mc.FieldMultipliers[identifier] = &multipler
			}
			
			d := metric(lo.Field(i).Interface(), ro.Field(i).Interface())
			total += d*multipler
		}
		return total
	}
	return f, nil
}

func (mc *MetricConfig) PointerMetric(lt, rt reflect.Type) (Metric, error) {
	f := func(l, r interface{})float64 {
		lp, rp := reflect.ValueOf(l), reflect.ValueOf(r)
		lv, rv := reflect.Indirect(lp), reflect.Indirect(rp)
		lti, rti := lv.Type(), rv.Type()
		metric, err := mc.GetMetric(lti, rti)
		if err != nil {
			panic(err)
		}
		return metric(lv.Interface(), rv.Interface())
	}
	return f, nil
}

func (mc *MetricConfig) InterfaceMetric() (Metric, error) {
	f := func(l, r interface{})float64{
		lt := reflect.TypeOf(l)
		rt := reflect.TypeOf(r)
		metric, err := mc.GetMetric(lt, rt)
		if err != nil {
			return mc.InterfaceCap
		} else {
			return math.Min(mc.InterfaceCap, metric(l, r))
		}
	}
	return f, nil
}
 
func (mc *MetricConfig) GetMetric(l, r reflect.Type) (Metric, error) {
	if l == nil && r == nil {
		return func(l,r interface{})float64{return 0}, nil
	}
	if l == mc.Hole || r == mc.Hole {
		return mc.DiscreteMetric(mc.HoleDistance), nil
	}
	if l != r {
		return nil, errors.New("unequal type, and neither is a hole")
	}
	switch l.Kind() {
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
	case reflect.Uintptr: return mc.DiscreteMetric(1.0), nil
	case reflect.Float32: return mc.NumberMetric(fromFloat32), nil
	case reflect.Float64: return mc.NumberMetric(fromFloat64), nil
	case reflect.Complex64: return mc.Complex64Metric, nil
	case reflect.Complex128: return mc.Complex128Metric, nil
	case reflect.Array: return mc.ArrayMetric(l, r)
	case reflect.Chan: return mc.ChanMetric(l, r)
	case reflect.Struct: return mc.StructMetric(l, r)
	case reflect.Ptr: return mc.PointerMetric(l, r)
	case reflect.Interface: return mc.InterfaceMetric() //return nil, errors.New(fmt.Sprintf("Interfaces can not directly be converted into metrics, consider wrapping in an object"))
	case reflect.String: return mc.StringMetric(), nil 
	case reflect.Slice: return mc.ArrayMetric(l, r)
	default: return nil, errors.New(fmt.Sprintf("undefined for type %s", l.String()))
	}
}