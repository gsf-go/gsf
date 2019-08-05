package component

import (
	"github.com/sf-go/gsf/src/gsc/serialization"
	"reflect"
	"strconv"
	"strings"
)

type Component struct {
	fields  map[string]reflect.Value
	version map[string]int
	record  map[string]interface{}

	objectId string
}

func NewComponent() *Component {
	return &Component{
		fields:  make(map[string]reflect.Value),
		version: make(map[string]int),
		record:  make(map[string]interface{}),
	}
}

func (component *Component) GetObjectId() string {
	return component.objectId
}

func (component *Component) Update(name string, value interface{}) bool {
	return true
}

func (component *Component) ToBinaryWriter(writer serialization.ISerializable) []byte {
	values := make([]interface{}, 0)

	for k, v := range component.record {
		values = append(values, k, v)
	}
	component.Clear()
	return writer.Serialize(values...)
}

func (component *Component) FromBinaryReader(reader serialization.IDeserializable) {
	values := reader.Deserialize()
	// Getter
	version := values[0].Interface().(string)
	splits := strings.Split(version, "_")
	length := len(splits) / 2
	remoteVersion := make(map[string]int)
	for i := 0; i < length; i += 2 {
		name := values[i].Interface().(string)
		ver, _ := strconv.Atoi(values[i+1].Interface().(string))
		remoteVersion[name] = ver

		if len(values) == 1 && component.version[name] > ver {
			component.record[name] = component.fields[name]
		}
	}

	// Setter
	values = values[1:]
	length = len(values) / 2
	for i := 0; i < length; i += 2 {
		name := values[i].Interface().(string)
		value := values[i+1].Interface()
		if remoteVersion[name] > component.version[name] && component.Update(name, values) {
			component.SetValue(name, value)
		}
	}
}

func (component *Component) Register(obj interface{}) {

	tv := reflect.ValueOf(obj).Elem()
	tk := reflect.TypeOf(obj).Elem()

	for i := 0; i < tv.NumField(); i++ {
		field := tv.Field(i)
		fieldType := field.Type()
		if fieldType == reflect.TypeOf(new(Component)) {
			continue
		}

		component.fields[tk.Field(i).Name] = field
	}
	component.objectId = reflect.TypeOf(obj).Elem().Name()
}

// 获得值
func (component *Component) GetValue(name string) interface{} {
	value, ok := component.fields[name]
	if ok {
		return value.Interface()
	}
	return nil
}

// 设置值
func (component *Component) SetValue(name string, value interface{}) {
	field, ok := component.fields[name]
	if ok {
		field.Set(reflect.ValueOf(value))
		component.version[name]++
		component.Record(name, value)
	}
}

// 记录值
func (component *Component) Record(name string, value interface{}) {
	component.record[name] = value
}

// 移除记录值
func (component *Component) Remove(name string) {
	delete(component.record, name)
}

// 清除所有记录值
func (component *Component) Clear() {
	for k := range component.record {
		component.Remove(k)
	}
}
