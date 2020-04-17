package main 

import (
	// "log"
	"go/types"
)

type GoToAxeConverter struct {
	Named map[*types.TypeName]Type
}

func (conv *GoToAxeConverter) GoVarToType(va *types.Var) Type {
	return conv.GoTypeToAxeType(va.Type())
}

func (conv *GoToAxeConverter) GoPointerToAxeKind(ptr *types.Pointer) *Kind {
	underType := conv.GoTypeToAxeType(ptr.Elem())
	return &Kind {
		name: "*",
		arguments: []Type{underType},
	}
}

func (conv *GoToAxeConverter) GoTupleToSliceOfAxeTypes(tup *types.Tuple) []Type {
	lentup := tup.Len()
	sliceTypes := make([]Type, lentup)
	for i := 0; i < lentup; i++ {
		va := tup.At(i)
		sliceTypes[i] = conv.GoVarToType(va)
	}
	return sliceTypes
}

func (conv *GoToAxeConverter) GoBasicToAxePrimative(basic *types.Basic) *Primative {
	return &Primative{basic.Name()}
}

func (conv *GoToAxeConverter) GoSignatureToAxeFunction(sig *types.Signature) *Function {
	return &Function {
		arguments: conv.GoTupleToSliceOfAxeTypes(sig.Params()),
		output: conv.GoTupleToSliceOfAxeTypes(sig.Results()),
	}
}

func (conv *GoToAxeConverter) GoArrayToAxeArray(array *types.Array) *Array {
	underlying := conv.GoTypeToAxeType(array.Elem())
	return &Array {
		size: int(array.Len()),
		typ: underlying,
	}
}

func (conv *GoToAxeConverter) GoSliceToAxeKind(slice *types.Slice) *Kind {
	underlying := conv.GoTypeToAxeType(slice.Elem())
	return &Kind{
		name: "slice",
		arguments: []Type{underlying},
	}
}

func (conv *GoToAxeConverter) GoMapToAxeKind(mp *types.Map) *Kind {
	from := conv.GoTypeToAxeType(mp.Key())
	to := conv.GoTypeToAxeType(mp.Elem())
	return &Kind {
		name: "map",
		arguments: []Type{from, to},
	}
}

func (conv *GoToAxeConverter) GoChanToAxeKind(ch *types.Chan) *Kind {
	underlying := conv.GoTypeToAxeType(ch.Elem())
	return &Kind {
		name: "chan",
		arguments: []Type{underlying},
	}
}

func (conv *GoToAxeConverter) GoStructToAxeStruct(strct *types.Struct) *Struct {
	nfields := strct.NumFields()
	fieldNames := make([]string, nfields)
	fields := make([]Type, nfields)
	for i := 0; i < nfields; i++ {
		va := strct.Field(i)
		fieldNames[i] = va.Name()
		fields[i] = conv.GoVarToType(va)
	}
	return &Struct {
		fieldNames: fieldNames,
		fields: fields,
	}
}

func (conv *GoToAxeConverter) GoNamedToAxeMethodHaver(named *types.Named) *MethodHaver {
	if val, ok := conv.Named[named.Obj()]; ok {
		return val.(*MethodHaver)
	}

	var mh MethodHaver
	conv.Named[named.Obj()] = &mh
	nmethods := named.NumMethods()
	methods := make(map[string]*Function, nmethods)
	for i := 0; i < nmethods; i++ {
		function := named.Method(i)
		functionSig := function.Type().(*types.Signature) // conversion promised by go spec
		axeFunction := conv.GoSignatureToAxeFunction(functionSig)
		methods[function.Name()] = axeFunction
	}	
	self := conv.GoTypeToAxeType(named.Underlying())
	mh.name = named.String()
	mh.self = self
	mh.methods = methods
	return &mh
}

func (conv *GoToAxeConverter) GoInterfaceToAxeInterface(gNtrf *types.Interface) *Interface  {
	var aNtrf Interface
	implements := make([]Type, gNtrf.NumEmbeddeds())
	for i := 0; i < gNtrf.NumEmbeddeds(); i++ {
		embedded := gNtrf.EmbeddedType(i)
		typ := conv.GoTypeToAxeType(embedded)
		implements[i] = typ 
	}
	
	methods := make(map[string]*Function, gNtrf.NumMethods())
	for i := 0; i < gNtrf.NumExplicitMethods(); i++ {
		function := gNtrf.ExplicitMethod(i)
		functionSig := function.Type().(*types.Signature)
		axeFunction := conv.GoSignatureToAxeFunction(functionSig)
		methods[function.Name()] = axeFunction
	}
	aNtrf.name = gNtrf.String()
	aNtrf.implements = implements
	aNtrf.methods = methods

	return &aNtrf
}

// converts .go types to actual types
func (conv *GoToAxeConverter) GoTypeToAxeType(gType types.Type) Type {
	// this switch is spelled out here, in the Types section:
	// https://github.com/golang/example/tree/master/gotypes
	switch gType.(type) {
	case *types.Basic: return conv.GoBasicToAxePrimative(gType.(*types.Basic))
	case *types.Pointer: return conv.GoPointerToAxeKind(gType.(*types.Pointer))
	case *types.Array: return conv.GoArrayToAxeArray(gType.(*types.Array))
	case *types.Slice: return conv.GoSliceToAxeKind(gType.(*types.Slice))
	case *types.Map: return conv.GoMapToAxeKind(gType.(*types.Map))
	case *types.Chan: return conv.GoChanToAxeKind(gType.(*types.Chan))
	case *types.Struct: return conv.GoStructToAxeStruct(gType.(*types.Struct))
	case *types.Signature: return conv.GoSignatureToAxeFunction(gType.(*types.Signature))
	case *types.Named: return conv.GoNamedToAxeMethodHaver(gType.(*types.Named))
	case *types.Interface: return conv.GoInterfaceToAxeInterface(gType.(*types.Interface))
	default: return &Hole{} // never reached
	}
}