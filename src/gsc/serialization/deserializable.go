package serialization

import (
	"github.com/sf-go/gsf/src/gsc/bytestream"
	"reflect"
)

type IDeserializable interface {
	Deserialize(bytes []byte, args ...interface{}) []reflect.Value
}

type Deserializable struct {
}

func NewDeserializable() *Deserializable {
	return &Deserializable{}
}

func (deserializable *Deserializable) Deserialize(bytes []byte, args ...interface{}) []reflect.Value {

	byteReader := bytestream.NewByteReader2(bytes)
	objects := make([]reflect.Value, 0)

	typeValue := uint8(0)
	byteReader.Read(&typeValue)
	kind := reflect.Kind(typeValue)
	length := deserializeValue(&kind, byteReader).Interface().(uint8)

	for i := uint8(0); i < length; i++ {
		value := DeserializeSingle(byteReader, args...)
		objects = append(objects, value)
	}

	return objects
}

func DeserializeSingle(byteReader *bytestream.ByteReader, args ...interface{}) reflect.Value {

	typeValue := uint8(0)
	byteReader.Read(&typeValue)
	kind := reflect.Kind(typeValue)

	// 反序列值类型
	if data := deserializeValue(&kind, byteReader); data != reflect.ValueOf(nil) {
		return data
	}

	// 反序列化引用类型
	if data := deserializeRef(&kind, byteReader); data != reflect.ValueOf(nil) {
		return data
	}

	// 反序列化切片类型
	if data := deserializeSlice(&kind, byteReader); data != reflect.ValueOf(nil) {
		return data
	}

	// 反序列化映射类型
	if data := deserializeMap(&kind, byteReader); data != reflect.ValueOf(nil) {
		return data
	}

	if data := deserializeStruct(&kind, byteReader, args...); data != reflect.ValueOf(nil) {
		return data
	}
	return reflect.ValueOf(nil)
}

func deserializeValue(
	kind *reflect.Kind,
	byteReader *bytestream.ByteReader) reflect.Value {

	if *kind == reflect.Struct {
		return reflect.ValueOf(nil)
	}

	generate, ok := GenerateVar[*kind]
	if !ok {
		return reflect.ValueOf(nil)
	}

	obj := generate(nil)
	byteReader.Read(obj)

	return reflect.ValueOf(obj).Elem()
}

func deserializeRef(
	kind *reflect.Kind,
	byteReader *bytestream.ByteReader) reflect.Value {

	if *kind != reflect.Ptr {
		return reflect.ValueOf(nil)
	}

	typeValue := uint8(0)
	byteReader.Read(&typeValue)
	*kind = reflect.Kind(typeValue)

	generate, ok := GenerateVar[*kind]
	if !ok {
		return reflect.ValueOf(nil)
	}

	obj := generate(nil)
	byteReader.Read(obj)

	return reflect.ValueOf(obj)
}

func deserializeSlice(
	kind *reflect.Kind,
	byteReader *bytestream.ByteReader,
) reflect.Value {

	if *kind != reflect.Slice {
		return reflect.ValueOf(nil)
	}

	typeValue := uint8(0)
	byteReader.Read(&typeValue)
	*kind = reflect.Kind(typeValue)

	if *kind == reflect.Ptr {

		byteReader.Read(&typeValue)
		*kind = reflect.Kind(typeValue)

		generate, ok := GenerateVar[*kind]
		if !ok {
			return reflect.ValueOf(nil)
		}

		valueLength := uint16(0)
		byteReader.Read(&valueLength)

		varType := reflect.SliceOf(KindPtrType[*kind])
		slice := reflect.MakeSlice(varType, int(valueLength), int(valueLength))

		for i := uint16(0); i < valueLength; i++ {
			obj := generate(nil)
			byteReader.Read(obj)
			slice.Index(int(i)).Set(reflect.ValueOf(obj))
		}
		return slice

	}

	generate, ok := GenerateVar[*kind]
	if !ok {
		return reflect.ValueOf(nil)
	}

	valueLength := uint16(0)
	byteReader.Read(&valueLength)

	varType := reflect.SliceOf(KindType[*kind])
	slice := reflect.MakeSlice(varType, int(valueLength), int(valueLength))

	for i := uint16(0); i < valueLength; i++ {
		obj := generate(nil)
		byteReader.Read(obj)
		slice.Index(int(i)).Set(reflect.ValueOf(obj).Elem())
	}
	return slice
}

func deserializeMap(kind *reflect.Kind, byteReader *bytestream.ByteReader) reflect.Value {
	if *kind != reflect.Map {
		return reflect.ValueOf(nil)
	}

	typeValue := uint8(0)
	byteReader.Read(&typeValue)

	keyKind := reflect.Kind(typeValue)
	keyGenerate, ok := GenerateVar[keyKind]
	if !ok {
		return reflect.ValueOf(nil)
	}

	byteReader.Read(&typeValue)
	valueKind := reflect.Kind(typeValue)

	if valueKind == reflect.Ptr {
		byteReader.Read(&typeValue)
		valueKind = reflect.Kind(typeValue)

		valueGenerate, ok := GenerateVar[valueKind]
		if !ok {
			return reflect.ValueOf(nil)
		}

		valueLength := uint16(0)
		byteReader.Read(&valueLength)

		varType := reflect.MapOf(KindType[keyKind], KindPtrType[valueKind])
		maps := reflect.MakeMap(varType)

		for i := uint16(0); i < valueLength; i++ {
			keyObj := keyGenerate(nil)
			valueObj := valueGenerate(nil)
			byteReader.Read(keyObj)
			byteReader.Read(valueObj)
			maps.SetMapIndex(reflect.ValueOf(keyObj).Elem(), reflect.ValueOf(valueObj).Elem())
		}
		return maps

	}

	valueGenerate, ok := GenerateVar[valueKind]
	if !ok {
		return reflect.ValueOf(nil)
	}

	valueLength := uint16(0)
	byteReader.Read(&valueLength)

	varType := reflect.MapOf(KindType[keyKind], KindType[valueKind])
	maps := reflect.MakeMap(varType)

	for i := uint16(0); i < valueLength; i++ {
		keyObj := keyGenerate(nil)
		valueObj := valueGenerate(nil)
		byteReader.Read(keyObj)
		byteReader.Read(valueObj)
		maps.SetMapIndex(reflect.ValueOf(keyObj).Elem(), reflect.ValueOf(valueObj).Elem())
	}
	return maps
}

func deserializeStruct(kind *reflect.Kind, byteReader *bytestream.ByteReader, args ...interface{}) reflect.Value {
	if *kind != reflect.Struct {
		return reflect.ValueOf(nil)
	}

	name := ""
	byteReader.Read(&name)

	valueGenerate, ok := GenerateVar[reflect.Struct]
	if !ok {
		return reflect.ValueOf(nil)
	}

	value := valueGenerate(append(append(make([]interface{}, 0), name), args...)...)
	packet := value.(ISerializablePacket)

	data := []byte(nil)
	byteReader.Read(&data)

	tmp := interface{}(nil)
	if len(args) == 0 {
		tmp = nil
	} else {
		tmp = args[0]
	}
	reader := NewEndianBinaryReader(data, tmp)
	packet.FromBinaryReader(reader)

	return reflect.ValueOf(packet)
}
