package serialization

import (
	"github.com/sf-go/gsf/src/gsc/bytestream"
	"reflect"
	"strconv"
	"testing"
)

func TestValueType(t *testing.T) {
	values := Show(t, 100, "hehe", nil)

	for _, v := range values {
		t.Log(v.Interface())
	}
}

func Show(t *testing.T, args ...interface{}) []reflect.Value {
	ser := new(Serializable)
	bytes := ser.Serialize(args...)
	des := NewDeserializable(bytestream.NewByteReader2(bytes))
	values := des.Deserialize()
	return values
}

func TestRefType(t *testing.T) {
	num := 100
	str := "hehe"

	values := Show(t, &num, &str)

	for _, v := range values {
		t.Log(v.Elem().Interface())
	}
}

func TestSliceRefType(t *testing.T) {
	array := make([]*bool, 2)
	b1 := false
	b2 := true
	array[0] = &b1
	array[1] = &b2

	values := Show(t, array)
	ret := values[0].Interface()
	for _, v := range ret.([]*bool) {
		t.Log(*v)
	}
}

func TestSliceValueType(t *testing.T) {
	array := make([]string, 2)
	array[0] = "xxxx"
	array[1] = "oooo"

	values := Show(t, array)
	ret := values[0].Interface()
	for _, v := range ret.([]string) {
		t.Log(v)
	}
}

func TestArrayValueType2(t *testing.T) {
	array := make([]float32, 2)
	array[0] = 3.1
	array[1] = 4.2

	values := Show(t, array)
	for _, v := range values[0].Interface().([]float32) {
		t.Log(v)
	}
}

func TestMapType(t *testing.T) {
	dict := make(map[string]int32)
	dict["xx"] = 50
	dict["oo"] = 50

	values := Show(t, dict)
	for k, v := range values[0].Interface().(map[string]int32) {
		t.Log(k + " " + strconv.Itoa(int(v)))
	}
}

type SerializablePacket struct {
	name string
	age  int
}

func NewSerializablePacket(name string, age int) *SerializablePacket {
	return &SerializablePacket{name: name, age: age}
}

func (serializablePacket *SerializablePacket) ToBinaryWriter(writer ISerializable) []byte {
	return writer.Serialize(serializablePacket.name, serializablePacket.age)
}

func (serializablePacket *SerializablePacket) FromBinaryReader(reader IDeserializable) {
	reader.Deserialize(&serializablePacket.name, &serializablePacket.age)
}

func TestStructType(t *testing.T) {
	GetPacketManagerInstance().AddPacket("SerializablePacket", func(args ...interface{}) ISerializablePacket {
		return NewSerializablePacket("", 0)
	})
	sut := NewSerializablePacket("Test", 100)
	values := Show(t, sut)
	packet := values[0].Interface().(*SerializablePacket)
	t.Log(packet.name + " " + strconv.Itoa(packet.age))
}

func TestStructType2(t *testing.T) {
	GetPacketManagerInstance().AddPacket("SerializablePacket", func(args ...interface{}) ISerializablePacket {
		return NewSerializablePacket("", 0)
	})

	sut := NewSerializablePacket("Test", 100)
	writer := NewSerializable()
	bytes := sut.ToBinaryWriter(writer)
	reader := NewDeserializable(bytestream.NewByteReader2(bytes))
	sut2 := NewSerializablePacket("", 0)
	sut2.FromBinaryReader(reader)
}

func TestSliceStruct(t *testing.T) {
	GetPacketManagerInstance().AddPacket("SerializablePacket", func(args ...interface{}) ISerializablePacket {
		return NewSerializablePacket("", 0)
	})
	sut := NewSerializablePacket("Test", 100)
	slice := []*SerializablePacket{
		sut,
		sut,
	}
	values := Show(t, slice)
	t.Log(values)
}

func TestMapStructType(t *testing.T) {

	GetPacketManagerInstance().AddPacket("SerializablePacket", func(args ...interface{}) ISerializablePacket {
		return NewSerializablePacket("", 0)
	})

	dict := make(map[string]*SerializablePacket)
	dict["xx"] = NewSerializablePacket("11", 100)
	dict["oo"] = NewSerializablePacket("22", 200)

	values := Show(t, dict)
	obj := values[0].Interface()
	for k, v := range obj.(map[string]*SerializablePacket) {
		t.Log(k + " " + string(v.name))
	}
}
