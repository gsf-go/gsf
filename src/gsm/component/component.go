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
	component := &Component{
		fields:  make(map[string]reflect.Value),
		version: make(map[string]int),
		record:  make(map[string]interface{}),
	}
	return component
}

func (component *Component) GetObjectId() string {
	return component.objectId
}

func (component *Component) Setter(name string, value interface{}) bool {
	return true
}

func (component *Component) Getter(version string) []interface{} {
	splits := strings.Split(version, "_")
	length := len(splits) / 2
	remoteVersion := make(map[string]int)
	tmp := make([]interface{}, 0)
	for i := 0; i < length; i += 2 {
		name := splits[i]
		ver, _ := strconv.Atoi(splits[i+1])
		remoteVersion[name] = ver

		if component.version[name] > ver {
			tmp = append(tmp, name+"_"+strconv.Itoa(component.version[name]))
			tmp = append(tmp, component.fields[name])
		}
	}
	return tmp
}

func (component *Component) Update() {

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
	length := len(values) / 2

	for i := 0; i < length; i += 2 {
		tmp := values[i].Interface().(string)
		splits := strings.Split(tmp, "_")
		name := splits[0]
		version, _ := strconv.Atoi(splits[1])
		value := values[i+1].Interface()
		if version > component.version[name] && component.Setter(name, values) {
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
