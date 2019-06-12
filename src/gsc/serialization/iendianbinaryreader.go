package serialization

import (
	"reflect"
)

type IEndianBinaryReader interface {
	Read(args ...interface{})
}

type EndianBinaryReader struct {
	deserializable IDeserializable
	buffer         []byte
}

func NewEndianBinaryReader(buffer []byte) *EndianBinaryReader {
	return &EndianBinaryReader{
		deserializable: NewDeserializable(),
		buffer:         buffer,
	}
}

func (reader *EndianBinaryReader) Read(args ...interface{}) {
	value := reader.deserializable.Deserialize(reader.buffer)
	for i, item := range args {
		reflect.ValueOf(item).Elem().Set(value[i])
	}
}
