package serialization

import (
	"reflect"
)

var GenerateVar = make(map[reflect.Kind]func(...interface{}) interface{})
var KindType = make(map[reflect.Kind]func(...interface{}) reflect.Type)
var KindPtrType = make(map[reflect.Kind]func(...interface{}) reflect.Type)

func init() {
	initVar()
	initType()
	initPtrType()
}

func initVar() {
	GenerateVar[reflect.Invalid] = func(params ...interface{}) interface{} {
		var invalid interface{}
		return &invalid
	}

	GenerateVar[reflect.Bool] = func(params ...interface{}) interface{} {
		b := new(bool)
		return b
	}

	GenerateVar[reflect.Int] = func(params ...interface{}) interface{} {
		i := new(int)
		return i
	}

	GenerateVar[reflect.Int8] = func(params ...interface{}) interface{} {
		i := new(int8)
		return i
	}

	GenerateVar[reflect.Int16] = func(params ...interface{}) interface{} {
		i := new(int16)
		return i
	}

	GenerateVar[reflect.Int32] = func(params ...interface{}) interface{} {
		i := new(int32)
		return i
	}

	GenerateVar[reflect.Int64] = func(params ...interface{}) interface{} {
		i := new(int64)
		return i
	}

	GenerateVar[reflect.Uint] = func(params ...interface{}) interface{} {
		i := new(uint)
		return i
	}

	GenerateVar[reflect.Uint8] = func(params ...interface{}) interface{} {
		i := new(uint8)
		return i
	}

	GenerateVar[reflect.Uint16] = func(params ...interface{}) interface{} {
		i := new(uint16)
		return i
	}

	GenerateVar[reflect.Uint32] = func(params ...interface{}) interface{} {
		i := new(uint32)
		return i
	}

	GenerateVar[reflect.Uint64] = func(params ...interface{}) interface{} {
		i := new(uint64)
		return i
	}

	GenerateVar[reflect.Float32] = func(params ...interface{}) interface{} {
		i := new(float32)
		return i
	}

	GenerateVar[reflect.Float64] = func(params ...interface{}) interface{} {
		i := new(float64)
		return i
	}

	GenerateVar[reflect.String] = func(params ...interface{}) interface{} {
		i := new(string)
		return i
	}

	GenerateVar[reflect.Struct] = func(params ...interface{}) interface{} {
		name := params[0].(string)
		packet := PacketManagerInstance.GetPacket(name)
		if packet != nil {
			return packet(params...)
		}
		return nil
	}
}

func initType() {
	KindType[reflect.Invalid] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(nil)
	}

	KindType[reflect.Bool] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(false)
	}

	KindType[reflect.Int] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(int(0))
	}

	KindType[reflect.Int8] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(int8(0))
	}

	KindType[reflect.Int16] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(int16(0))
	}

	KindType[reflect.Int32] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(int32(0))
	}

	KindType[reflect.Int64] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(int64(0))
	}

	KindType[reflect.Uint] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(uint(0))
	}

	KindType[reflect.Uint8] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(uint8(0))
	}

	KindType[reflect.Uint16] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(uint16(0))
	}

	KindType[reflect.Uint32] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(uint32(0))
	}

	KindType[reflect.Uint64] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(uint64(0))
	}

	KindType[reflect.Float32] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(float32(0))
	}

	KindType[reflect.Float64] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(float64(0))
	}

	KindType[reflect.String] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf("")
	}

	KindType[reflect.Struct] = func(params ...interface{}) reflect.Type {
		name := params[0].(string)
		packet := PacketManagerInstance.GetPacket(name)
		if packet != nil {
			return reflect.TypeOf(packet(params...))
		}
		return nil
	}
}

func initPtrType() {

	Invalid := new(interface{})
	KindPtrType[reflect.Invalid] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(&Invalid)
	}

	Bool := false
	KindPtrType[reflect.Bool] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(&Bool)
	}

	Int := int(0)
	KindPtrType[reflect.Int] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(&Int)
	}

	Int8 := int8(0)
	KindPtrType[reflect.Int8] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(&Int8)
	}

	Int16 := int16(0)
	KindPtrType[reflect.Int16] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(&Int16)
	}

	Int32 := int32(0)
	KindPtrType[reflect.Int32] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(&Int32)
	}

	Int64 := int64(0)
	KindPtrType[reflect.Int64] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(&Int64)
	}

	Uint := uint(0)
	KindPtrType[reflect.Uint] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(&Uint)
	}

	Uint8 := uint8(0)
	KindPtrType[reflect.Uint8] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(&Uint8)
	}

	Uint16 := uint16(0)
	KindPtrType[reflect.Uint16] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(&Uint16)
	}

	Uint32 := uint32(0)
	KindPtrType[reflect.Uint32] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(&Uint32)
	}

	Uint64 := uint64(0)
	KindPtrType[reflect.Uint64] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(&Uint64)
	}

	Float32 := float32(0)
	KindPtrType[reflect.Float32] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(&Float32)
	}

	Float64 := float64(0)
	KindPtrType[reflect.Float64] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(&Float64)
	}

	String := ""
	KindPtrType[reflect.String] = func(params ...interface{}) reflect.Type {
		return reflect.TypeOf(&String)
	}

	KindPtrType[reflect.Struct] = func(params ...interface{}) reflect.Type {
		name := params[0].(string)
		packet := PacketManagerInstance.GetPacket(name)
		if packet != nil {
			return reflect.TypeOf(packet(params...))
		}
		return nil
	}
}
