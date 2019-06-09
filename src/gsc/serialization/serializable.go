package serialization

import (
	"gsc/bytestream"
	"reflect"
)

type ISerializable interface {
	Serialize(args ...interface{}) []byte
}

type Serializable struct {
}

func (serializable *Serializable) Serialize(args ...interface{}) []byte {

	buffer := make([]byte, 0)
	bytes := serializeValue(uint16(len(args)))
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
	}

	return buffer
}

func serializeValue(value interface{}) []byte {

	king := reflect.ValueOf(value).Kind()
	if _, ok := GenerateVar[king]; !ok {
		return nil
	}

	byteWriter := bytestream.NewByteWriter3()
	compositeType := NewCompositeType(make([]byte, 0))

	compositeType.Append(uint8(king))
	byteWriter.Write(compositeType.GetType())
	byteWriter.Write(value)
	return byteWriter.ToBytes()
}

func serializeRef(value interface{}) []byte {

	valueType := reflect.TypeOf(value)
	if valueType.Kind() != reflect.Ptr {
		return nil
	}

	byteWriter := bytestream.NewByteWriter3()
	compositeType := NewCompositeType(make([]byte, 0))

	compositeType.Append(uint8(reflect.Ptr))
	compositeType.Append(uint8(valueType.Elem().Kind()))
	byteWriter.Write(compositeType.GetType())
	byteWriter.Write(value)
	return byteWriter.ToBytes()
}

func serializeSlice(value interface{}) []byte {

	varType := reflect.TypeOf(value)
	if varType.Kind() != reflect.Slice {
		return nil
	}

	byteWriter := bytestream.NewByteWriter3()
	compositeType := NewCompositeType(make([]byte, 0))

	compositeType.Append(uint8(reflect.Slice))
	kind := varType.Elem().Kind()
	compositeType.Append(uint8(kind))

	if kind == reflect.Ptr {
		kind = varType.Elem().Elem().Kind()
		compositeType.Append(uint8(kind))
	}

	byteWriter.Write(compositeType.GetType())
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
	compositeType := NewCompositeType(make([]byte, 0))
	compositeType.Append(uint8(reflect.Map))

	keyKind := varType.Key().Kind()
	compositeType.Append(uint8(keyKind))
	valueKind := varType.Elem().Kind()
	compositeType.Append(uint8(valueKind))

	if valueKind == reflect.Ptr {
		kind := varType.Elem().Elem().Kind()
		compositeType.Append(uint8(kind))
	}

	byteWriter.Write(compositeType.GetType())
	valueType := reflect.ValueOf(value)
	length := uint16(valueType.Len())
	byteWriter.Write(length)

	for _, key := range valueType.MapKeys() {
		byteWriter.Write(key.Interface())
		byteWriter.Write(valueType.MapIndex(key).Interface())
	}

	return byteWriter.ToBytes()
}
