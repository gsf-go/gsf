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

func TestSize(t *testing.T) {
	bw := NewByteWriter2(make([]byte, 0))
	num1 := int32(100)
	num2 := int32(200)

	bw.Write(num1)
	bw.Write(num2)

	num1 = 0
	num2 = 0

	br := NewByteReader2(bw.ToBytes())
	br.Read(&num1)
	br.Read(&num1)

	br.Shift(-4)
	pos := br.GetPosition()
	t.Logf("%d", pos)
}
