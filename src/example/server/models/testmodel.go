package models

import "github.com/sf-go/gsf/src/gsc/serialization"

type TestModel struct {
	Name string
	Age  int
}

func NewTestModel() *TestModel {
	return &TestModel{}
}

func (testModel *TestModel) ToBinaryWriter(writer serialization.IEndianBinaryWriter) {
	writer.Write(testModel.Name, testModel.Age)
}

func (testModel *TestModel) FromBinaryReader(reader serialization.IEndianBinaryReader) {
	reader.Read(&testModel.Name, &testModel.Age)
}
