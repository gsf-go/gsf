package bytestream

import (
	"testing"
)

func TestBytes(t *testing.T) {
	bw := NewByteWriter2(make([]byte, 0))
	num := int32(100)
	str := "test"
	flt := 2.0
	bw.Write(num)
	bw.Write(str)
	bw.Write(flt)

	num = int32(0)
	str = ""
	flt = 0.0

	br := NewByteReader2(bw.ToBytes())
	br.Read(&num)
	br.Read(&str)

	br.Shift(-6)
	str2 := ""
	br.Read(&str2)
	br.Read(&flt)
}
