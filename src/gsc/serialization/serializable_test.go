package serialization

import (
	"testing"
)

func TestValueType(t *testing.T) {
	Show(t, 100, "hehe", nil)
}

func Show(t *testing.T, args ...interface{}) {
	ser := new(Serializable)
	deser := new(Deserializable)

	bytes := ser.Serialize(args)
	value := deser.Deserialize(bytes)
	t.Log(value)
}

func TestRefType(t *testing.T) {
	num := 100
	str := "hehe"

	Show(t, &num, &str)
}

func TestSliceRefType(t *testing.T) {
	array := make([]*bool, 2)
	b1 := false
	b2 := true
	array[0] = &b1
	array[1] = &b2

	Show(t, array)
}

func TestSliceValueType(t *testing.T) {
	array := make([]string, 2)
	array[0] = "xxxx"
	array[1] = "oooo"

	Show(t, array)
}

func TestArrayValueType2(t *testing.T) {
	array := make([]float32, 2)
	array[0] = 3.1
	array[1] = 4.2

	Show(t, array)
}

func TestMapType(t *testing.T) {
	dict := make(map[string]int32)
	dict["xx"] = 50
	dict["oo"] = 50

	Show(t, dict)
}
