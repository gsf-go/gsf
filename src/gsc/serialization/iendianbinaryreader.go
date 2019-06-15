package serialization

import (
	"reflect"
)

type IEndianBinaryReader interface {
	Read(args ...interface{})
	ReadValues() []interface{}
}

type EndianBinaryReader struct {
	deserializable IDeserializable
	buffer         []byte
	extension      interface{}
}

func NewEndianBinaryReader(buffer []byte, extension interface{}) *EndianBinaryReader {
	return &EndianBinaryReader{
		deserializable: NewDeserializable(),
		buffer:         buffer,
		extension:      extension,
	}
}

func (reader *EndianBinaryReader) Read(args ...interface{}) {
	value := reader.deserializable.Deserialize(reader.buffer)
	for i, item := range args {
		reflect.ValueOf(item).Elem().Set(value[i])
	}
}

func (reader *EndianBinaryReader) ReadValues() []interface{} {
	values := reader.deserializable.Deserialize(reader.buffer, reader.extension)
	result := make([]interface{}, 0)
	for _, item := range values {
		result = append(result, item.Interface())
	}
	return result
}
