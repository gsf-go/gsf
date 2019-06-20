package models

import "github.com/sf-go/gsf/src/gsc/serialization"

type TestModel struct {
	Name string
	Age  int
}

func NewTestModel() *TestModel {
	return &TestModel{}
}

func (testModel *TestModel) ToBinaryWriter(writer serialization.ISerializable) []byte {
	return writer.Serialize(testModel.Name, testModel.Age)
}

func (testModel *TestModel) FromBinaryReader(reader serialization.IDeserializable) {
	reader.Deserialize(&testModel.Name, &testModel.Age)
}
