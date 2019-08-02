package serialization

import (
	"github.com/sf-go/gsf/src/gsc/bytestream"
	"reflect"
)

type ISerializable interface {
	Serialize(args ...interface{}) []byte
	SerializeSingle(obj interface{}) []byte
}

type Serializable struct {
}

func NewSerializable() *Serializable {
	return &Serializable{}
}

func (serializable *Serializable) Serialize(args ...interface{}) []byte {

	buffer := make([]byte, 0)
	bytes := serializable.serializeValue(uint8(len(args)))
	buffer = append(buffer, bytes...)

	for _, item := range args {

		bytes = serializable.serializeValue(item)
		if bytes != nil {
			buffer = append(buffer, bytes...)
			continue
		}

		bytes = serializable.serializeRef(item)
		if bytes != nil {
			buffer = append(buffer, bytes...)
			continue
		}

		bytes = serializable.serializeSlice(item)
		if bytes != nil {
			buffer = append(buffer, bytes...)
			continue
		}

		bytes = serializable.serializeMap(item)
		if bytes != nil {
			buffer = append(buffer, bytes...)
			continue
		}

		bytes = serializable.serializeStruct(item)
		if bytes != nil {
			buffer = append(buffer, bytes...)
			continue
		}
	}

	return buffer
}

func (serializable *Serializable) SerializeSingle(obj interface{}) []byte {

	bytes := serializable.serializeValue(obj)
	if bytes != nil {
		return bytes
	}

	bytes = serializable.serializeRef(obj)
	if bytes != nil {
		return bytes
	}

	bytes = serializable.serializeSlice(obj)
	if bytes != nil {
		return bytes
	}

	bytes = serializable.serializeMap(obj)
	if bytes != nil {
		return bytes
	}

	bytes = serializable.serializeStruct(obj)
	if bytes != nil {
		return bytes
	}

	return bytes
}

func (serializable *Serializable) serializeValue(value interface{}) []byte {

	kind := reflect.ValueOf(value).Kind()
	if _, ok := GenerateVar[kind]; !ok {
		return nil
	}

	byteWriter := bytestream.NewByteWriter3()
	byteWriter.Write(uint8(kind))
	byteWriter.Write(value)
	return byteWriter.ToBytes()
}

func (serializable *Serializable) serializeRef(value interface{}) []byte {

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

func (serializable *Serializable) serializeSlice(value interface{}) []byte {

	varType := reflect.TypeOf(value)
	if varType.Kind() != reflect.Slice {
		return nil
	}

	byteWriter := bytestream.NewByteWriter3()
	byteWriter.Write(uint8(reflect.Slice))

	kind := varType.Elem().Kind()

	if kind == reflect.Ptr {
		byteWriter.Write(uint8(reflect.Ptr))
		element := varType.Elem().Elem()
		kind = element.Kind()
		if kind == reflect.Struct {
			byteWriter.Write(uint8(kind))
			byteWriter.Write(element.Name())
		} else {
			byteWriter.Write(uint8(kind))
		}
	} else {
		byteWriter.Write(uint8(varType.Elem().Kind()))
	}

	valueType := reflect.ValueOf(value)
	length := uint16(valueType.Len())
	byteWriter.Write(length)

	for i := 0; i < int(length); i++ {
		item := valueType.Index(i)
		v := item.Interface()
		kind := item.Kind()
		if kind == reflect.Ptr && item.Elem().Kind() == reflect.Struct {
			byteWriter.Write(serializable.serializeStruct(v))
		} else {
			byteWriter.Write(v)
		}
	}

	return byteWriter.ToBytes()
}

func (serializable *Serializable) serializeMap(value interface{}) []byte {

	mapType := reflect.TypeOf(value)
	if mapType.Kind() != reflect.Map {
		return nil
	}

	byteWriter := bytestream.NewByteWriter3()
	byteWriter.Write(uint8(reflect.Map))

	valueKind := mapType.Elem().Kind()
	kind := reflect.Invalid
	if valueKind == reflect.Ptr {
		byteWriter.Write(uint8(reflect.Ptr))
		valueType := mapType.Elem().Elem()
		kind = valueType.Kind()

		byteWriter.Write(uint8(kind))
		if kind == reflect.Struct {
			byteWriter.Write(valueType.Name())
		}
	} else {
		kind = mapType.Elem().Kind()
		byteWriter.Write(uint8(kind))
		if kind == reflect.Struct {
			byteWriter.Write(mapType.Elem().Name())
		}
	}

	keyKind := mapType.Key().Kind()
	byteWriter.Write(uint8(keyKind))

	valueType := reflect.ValueOf(value)
	length := uint16(valueType.Len())
	byteWriter.Write(length)

	for _, key := range valueType.MapKeys() {
		byteWriter.Write(key.Interface())
		if kind == reflect.Struct {
			byteWriter.Write(serializable.serializeStruct(valueType.MapIndex(key).Interface()))
		} else {
			byteWriter.Write(valueType.MapIndex(key).Interface())
		}
	}

	return byteWriter.ToBytes()
}

func (serializable *Serializable) serializeStruct(value interface{}) []byte {
	varType := reflect.TypeOf(value)
	if varType.Kind() != reflect.Ptr || varType.Elem().Kind() != reflect.Struct {
		return nil

	}
	serializablePacket := value.(ISerializablePacket)
	if serializablePacket == nil {
		return nil
	}

	byteWriter := bytestream.NewByteWriter3()
	bytes := serializablePacket.ToBinaryWriter(serializable)

	byteWriter.Write(uint8(reflect.Struct))
	byteWriter.Write(varType.Elem().Name())
	byteWriter.Write(bytes)

	return byteWriter.ToBytes()
}
