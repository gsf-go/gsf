package property

import (
	"github.com/gsf/gsf/src/gsc/serialization"
	"reflect"
	"sync"
)

type ISetter interface {
	SetValue(name string, value interface{})
}

type IGetter interface {
	GetValue(name string) interface{}
}

type Property struct {
	fields *sync.Map
	record *sync.Map
}

func NewProperty() *Property {
	return &Property{
		fields: new(sync.Map),
		record: new(sync.Map),
	}
}

func (property *Property) Register(obj interface{}) {

	tv := reflect.ValueOf(obj).Elem()
	tk := reflect.TypeOf(obj).Elem()

	for i := 0; i < tv.NumField(); i++ {
		field := tv.Field(i)
		fieldType := field.Type()
		if fieldType == reflect.TypeOf(new(Property)) {
			continue
		}

		property.fields.Store(tk.Field(i).Name, field)
	}
}

func (property *Property) GetValue(name string) interface{} {
	value, ok := property.fields.Load(name)
	if ok {
		return value.(reflect.Value).Interface()
	}
	return nil
}

func (property *Property) SetValue(name string, value interface{}) {
	field, ok := property.fields.Load(name)
	if ok {
		fieldValue := field.(reflect.Value)
		fieldValue.Set(reflect.ValueOf(value))
		property.Record(name, value)
	}
}

func (property *Property) Record(name string, value interface{}) {
	property.record.Store(name, value)
}

func (property *Property) Remove(name string) {
	property.record.Delete(name)
}

func (property *Property) Clear() {
	property.record.Range(func(key, value interface{}) bool {
		property.Remove(key.(string))
		return true
	})
}

func (property *Property) ToBinaryWriter(writer serialization.IEndianBinaryWriter) {
	values := make([]interface{}, 0)
	property.record.Range(func(key, value interface{}) bool {
		values = append(values, key, value)
		return true
	})
	property.Clear()
	writer.Write(values...)
}

func (property *Property) FromBinaryReader(reader serialization.IEndianBinaryReader) {
	values := reader.ReadValues()
	length := len(values) / 2
	for i := 0; i < length; i++ {
		name := values[i*2].(string)
		value := values[i*2+1]
		property.SetValue(name, value)
	}
}
