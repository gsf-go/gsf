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

func (byteReader *ByteReader) IsEOF() bool {
	return byteReader.reader.Len() == 0
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

		dataLength := int(length) - binary.Size(uint16(0))
		bytes := make([]byte, dataLength)

		err = binary.Read(byteReader.reader, byteReader.Order, bytes)
		if err != nil {
			panic(err)
		}
		*data = string(bytes)
	case *interface{}:
		return

	default:
		err := binary.Read(byteReader.reader, byteReader.Order, data)
		if err != nil {
			panic(err)
		}
	}
}