package bytestream

import (
	"bytes"
	"encoding/binary"
	"io"
)

type ByteReader struct {
	Order  binary.ByteOrder
	reader *bytes.Reader
}

func NewByteReader(buffer []byte, order binary.ByteOrder) *ByteReader {
	return &ByteReader{
		Order:  order,
		reader: bytes.NewReader(buffer),
	}
}

func NewByteReader2(buffer []byte) *ByteReader {
	return NewByteReader(buffer, binary.LittleEndian)
}

func (byteReader *ByteReader) Shift(offset int64) {
	_, err := byteReader.reader.Seek(offset, io.SeekCurrent)
	if err != nil {
		panic(err)
	}
}

func (byteReader *ByteReader) GetPosition() int {
	return byteReader.reader.Len()
}

func (byteReader *ByteReader) Read(data interface{}) {
	switch data := data.(type) {
	case *int:
		var v int32 = 0
		err := binary.Read(byteReader.reader, byteReader.Order, &v)
		if err != nil {
			panic(err)
		}

		*data = int(v)
	case *uint:
		v := uint32(0)
		err := binary.Read(byteReader.reader, byteReader.Order, &v)
		if err != nil {
			panic(err)
		}

		*data = uint(v)
	case *string:
		length := uint16(0)
		err := binary.Read(byteReader.reader, byteReader.Order, &length)
		if err != nil {
			panic(err)
		}

		bytes := make([]byte, length)

		err = binary.Read(byteReader.reader, byteReader.Order, bytes)
		if err != nil {
			panic(err)
		}

		*data = string(bytes)
	case *interface{}:
		return

	case *[]byte:
		length := uint16(0)
		err := binary.Read(byteReader.reader, byteReader.Order, &length)
		if err != nil {
			panic(err)
		}

		bytes := make([]byte, length)

		err = binary.Read(byteReader.reader, byteReader.Order, bytes)
		if err != nil {
			panic(err)
		}

		*data = bytes
	default:
		err := binary.Read(byteReader.reader, byteReader.Order, data)
		if err != nil {
			panic(err)
		}
	}
}
