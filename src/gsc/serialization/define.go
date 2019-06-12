package serialization

import "reflect"

var GenerateVar = make(map[reflect.Kind]func(interface{}) interface{})
var KindType = make(map[reflect.Kind]reflect.Type)
var KindPtrType = make(map[reflect.Kind]reflect.Type)

func init() {
	initVar()
	initType()
	initPtrType()
}

func initVar() {
	GenerateVar[reflect.Invalid] = func(param interface{}) interface{} {
		var invalid interface{}
		return &invalid
	}

	GenerateVar[reflect.Bool] = func(param interface{}) interface{} {
		b := new(bool)
		return b
	}

	GenerateVar[reflect.Int] = func(param interface{}) interface{} {
		i := new(int)
		return i
	}

	GenerateVar[reflect.Int8] = func(param interface{}) interface{} {
		i := new(int8)
		return i
	}

	GenerateVar[reflect.Int16] = func(param interface{}) interface{} {
		i := new(int16)
		return i
	}

	GenerateVar[reflect.Int32] = func(param interface{}) interface{} {
		i := new(int32)
		return i
	}

	GenerateVar[reflect.Int64] = func(param interface{}) interface{} {
		i := new(int64)
		return i
	}

	GenerateVar[reflect.Uint] = func(param interface{}) interface{} {
		i := new(uint)
		return i
	}

	GenerateVar[reflect.Uint8] = func(param interface{}) interface{} {
		i := new(uint8)
		return i
	}

	GenerateVar[reflect.Uint16] = func(param interface{}) interface{} {
		i := new(uint16)
		return i
	}

	GenerateVar[reflect.Uint32] = func(param interface{}) interface{} {
		i := new(uint32)
		return i
	}

	GenerateVar[reflect.Uint64] = func(param interface{}) interface{} {
		i := new(uint64)
		return i
	}

	GenerateVar[reflect.Float32] = func(param interface{}) interface{} {
		i := new(float32)
		return i
	}

	GenerateVar[reflect.Float64] = func(param interface{}) interface{} {
		i := new(float64)
		return i
	}

	GenerateVar[reflect.String] = func(param interface{}) interface{} {
		i := new(string)
		return i
	}

	GenerateVar[reflect.Struct] = func(param interface{}) interface{} {
		name := param.(string)
		packet := GetPacketManagerInstance().GetPacket(name)
		if packet != nil {
			return packet()
		}
		return nil
	}
}

func initType() {
	KindType[reflect.Invalid] = reflect.TypeOf(nil)
	KindType[reflect.Bool] = reflect.TypeOf(false)
	KindType[reflect.Int] = reflect.TypeOf(int(0))
	KindType[reflect.Int8] = reflect.TypeOf(int8(0))
	KindType[reflect.Int16] = reflect.TypeOf(int16(0))
	KindType[reflect.Int32] = reflect.TypeOf(int32(0))
	KindType[reflect.Int64] = reflect.TypeOf(int64(0))
	KindType[reflect.Uint] = reflect.TypeOf(uint(0))
	KindType[reflect.Uint8] = reflect.TypeOf(uint8(0))
	KindType[reflect.Uint16] = reflect.TypeOf(uint16(0))
	KindType[reflect.Uint32] = reflect.TypeOf(uint32(0))
	KindType[reflect.Uint64] = reflect.TypeOf(uint64(0))
	KindType[reflect.Float32] = reflect.TypeOf(float32(0))
	KindType[reflect.Float64] = reflect.TypeOf(float64(0))
	KindType[reflect.String] = reflect.TypeOf("")
}

func initPtrType() {
	Invalid := new(interface{})
	KindPtrType[reflect.Invalid] = reflect.TypeOf(&Invalid)
	Bool := false
	KindPtrType[reflect.Bool] = reflect.TypeOf(&Bool)
	Int := int(0)
	KindPtrType[reflect.Int] = reflect.TypeOf(&Int)
	Int8 := int8(0)
	KindPtrType[reflect.Int8] = reflect.TypeOf(&Int8)
	Int16 := int16(0)
	KindPtrType[reflect.Int16] = reflect.TypeOf(&Int16)
	Int32 := int32(0)
	KindPtrType[reflect.Int32] = reflect.TypeOf(&Int32)
	Int64 := int64(0)
	KindPtrType[reflect.Int64] = reflect.TypeOf(&Int64)
	Uint := uint(0)
	KindPtrType[reflect.Uint] = reflect.TypeOf(&Uint)
	Uint8 := uint8(0)
	KindPtrType[reflect.Uint8] = reflect.TypeOf(&Uint8)
	Uint16 := uint16(0)
	KindPtrType[reflect.Uint16] = reflect.TypeOf(&Uint16)
	Uint32 := uint32(0)
	KindPtrType[reflect.Uint32] = reflect.TypeOf(&Uint32)
	Uint64 := uint64(0)
	KindPtrType[reflect.Uint64] = reflect.TypeOf(&Uint64)
	Float32 := float32(0)
	KindPtrType[reflect.Float32] = reflect.TypeOf(&Float32)
	Float64 := float64(0)
	KindPtrType[reflect.Float64] = reflect.TypeOf(&Float64)
	String := ""
	KindPtrType[reflect.String] = reflect.TypeOf(&String)
}
