package bytestream

import (
	"bytes"
	"encoding/binary"
)

type ByteWriter struct {
	Order  binary.ByteOrder
	buffer *bytes.Buffer
}

func NewByteWriter(buffer []byte, order binary.ByteOrder) *ByteWriter {
	return &ByteWriter{
		Order:  order,
		buffer: bytes.NewBuffer(buffer),
	}
}

func NewByteWriter2(buffer []byte) *ByteWriter {
	return NewByteWriter(buffer, binary.LittleEndian)
}

func NewByteWriter3() *ByteWriter {
	return &ByteWriter{
		Order:  binary.LittleEndian,
		buffer: new(bytes.Buffer),
	}
}

func (byteWriter *ByteWriter) Write(data interface{}) {
	switch data := data.(type) {
	case *int:
		err := binary.Write(byteWriter.buffer, byteWriter.Order, int32(*data))
		if err != nil {
			panic(err)
		}
	case *uint:
		err := binary.Write(byteWriter.buffer, byteWriter.Order, uint32(*data))
		if err != nil {
			panic(err)
		}
	case int:
		err := binary.Write(byteWriter.buffer, byteWriter.Order, int32(data))
		if err != nil {
			panic(err)
		}
	case uint:
		err := binary.Write(byteWriter.buffer, byteWriter.Order, uint32(data))
		if err != nil {
			panic(err)
		}
	case string:
		length := uint16(len(data))
		err := binary.Write(byteWriter.buffer, byteWriter.Order, length)
		if err != nil {
			panic(err)
		}

		buffer := []byte(data)
		err = binary.Write(byteWriter.buffer, byteWriter.Order, buffer)
		if err != nil {
			panic(err)
		}
	case *string:
		length := uint16(len(*data))
		err := binary.Write(byteWriter.buffer, byteWriter.Order, length)
		if err != nil {
			panic(err)
		}

		buffer := []byte(*data)
		err = binary.Write(byteWriter.buffer, byteWriter.Order, buffer)
		if err != nil {
			panic(err)
		}
	case nil:
		return
	case []byte:
		length := uint16(len(data))
		err := binary.Write(byteWriter.buffer, byteWriter.Order, length)
		if err != nil {
			panic(err)
		}

		buffer := []byte(data)
		err = binary.Write(byteWriter.buffer, byteWriter.Order, buffer)
		if err != nil {
			panic(err)
		}
	default:
		err := binary.Write(byteWriter.buffer, byteWriter.Order, data)
		if err != nil {
			panic(err)
		}
	}
}

func (byteWriter *ByteWriter) ToBytes() []byte {
	return byteWriter.buffer.Bytes()
}
