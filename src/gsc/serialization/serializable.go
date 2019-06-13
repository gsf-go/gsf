package serialization

import (
	"github.com/gsf/gsf/src/gsc/bytestream"
	"reflect"
)

type ISerializable interface {
	Serialize(args ...interface{}) []byte
}

type Serializable struct {
}

func NewSerializable() *Serializable {
	return &Serializable{}
}

func (serializable *Serializable) Serialize(args ...interface{}) []byte {

	buffer := make([]byte, 0)
	bytes := serializeValue(uint8(len(args)))
	buffer = append(buffer, bytes...)

	for _, item := range args {

		bytes = serializeValue(item)
		if bytes != nil {
			buffer = append(buffer, bytes...)
			continue
		}

		bytes = serializeRef(item)
		if bytes != nil {
			buffer = append(buffer, bytes...)
			continue
		}

		bytes = serializeSlice(item)
		if bytes != nil {
			buffer = append(buffer, bytes...)
			continue
		}

		bytes = serializeMap(item)
		if bytes != nil {
			buffer = append(buffer, bytes...)
			continue
		}

		bytes = serializeStruct(item)
		if bytes != nil {
			buffer = append(buffer, bytes...)
			continue
		}
	}

	return buffer
}

func serializeValue(value interface{}) []byte {

	king := reflect.ValueOf(value).Kind()
	if _, ok := GenerateVar[king]; !ok {
		return nil
	}

	byteWriter := bytestream.NewByteWriter3()
	byteWriter.Write(uint8(king))
	byteWriter.Write(value)
	return byteWriter.ToBytes()
}

func serializeRef(value interface{}) []byte {

	valueType := reflect.TypeOf(value)
	if valueType.Kind() != reflect.Ptr || valueType.Elem().Kind() == reflect.Struct {
		return nil
	}

	byteWriter := bytestream.NewByteWriter3()
	byteWriter.Write(uint8(reflect.Ptr))
	byteWriter.Write(uint8(valueType.Elem().Kind()))
	byteWriter.Write(value)
	return byteWriter.ToBytes()
}

func serializeSlice(value interface{}) []byte {

	varType := reflect.TypeOf(value)
	if varType.Kind() != reflect.Slice {
		return nil
	}

	byteWriter := bytestream.NewByteWriter3()
	byteWriter.Write(uint8(reflect.Slice))

	kind := varType.Elem().Kind()

	if kind == reflect.Ptr {
		byteWriter.Write(uint8(reflect.Ptr))
		byteWriter.Write(uint8(varType.Elem().Elem().Kind()))
	} else {
		byteWriter.Write(uint8(varType.Elem().Kind()))
	}

	valueType := reflect.ValueOf(value)
	length := uint16(valueType.Len())
	byteWriter.Write(length)

	for i := 0; i < int(length); i++ {
		v := valueType.Index(i).Interface()
		byteWriter.Write(v)
	}

	return byteWriter.ToBytes()
}

func serializeMap(value interface{}) []byte {

	varType := reflect.TypeOf(value)
	if varType.Kind() != reflect.Map {
		return nil
	}

	byteWriter := bytestream.NewByteWriter3()
	byteWriter.Write(uint8(reflect.Map))

	keyKind := varType.Key().Kind()
	byteWriter.Write(uint8(keyKind))
	valueKind := varType.Elem().Kind()

	if valueKind == reflect.Ptr {
		byteWriter.Write(uint8(varType.Elem().Elem().Kind()))
	} else {
		byteWriter.Write(uint8(varType.Elem().Kind()))
	}

	valueType := reflect.ValueOf(value)
	length := uint16(valueType.Len())
	byteWriter.Write(length)

	for _, key := range valueType.MapKeys() {
		byteWriter.Write(key.Interface())
		byteWriter.Write(valueType.MapIndex(key).Interface())
	}

	return byteWriter.ToBytes()
}

func serializeStruct(value interface{}) []byte {
	varType := reflect.TypeOf(value)
	if varType.Kind() != reflect.Ptr || varType.Elem().Kind() != reflect.Struct {
		return nil
	}

	serializablePacket := value.(ISerializablePacket)
	if serializablePacket == nil {
		return nil
	}

	byteWriter := bytestream.NewByteWriter3()

	writer := NewEndianBinaryWriter()
	serializablePacket.ToBinaryWriter(writer)
	bytes := writer.ToBytes()

	byteWriter.Write(uint8(reflect.Struct))
	byteWriter.Write(reflect.TypeOf(value).Elem().Name())
	byteWriter.Write(bytes)

	return byteWriter.ToBytes()
}
