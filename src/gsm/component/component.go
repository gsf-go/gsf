package component

import (
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

func (component *Component) Synchronize() (string, []interface{}) {
	values := make([]interface{}, 0)

	for name, value := range component.record {
		values = append(values, name+"_"+strconv.Itoa(component.version[name]))
		values = append(values, value)
	}
	component.Clear()
	return "Get_" + component.objectId, values
}

func (component *Component) GetObjectId() string {
	return component.objectId
}

func (component *Component) Verify(name string, value interface{}) bool {
	return true
}

func (component *Component) GetterCallback(version string) []interface{} {
	splits := strings.Split(version, "_")
	length := len(splits)
	tmp := make([]interface{}, 0)
	for i := 0; i < length; i += 2 {
		name := splits[i]
		ver, _ := strconv.Atoi(splits[i+1])
		if component.version[name] > ver {
			tmp = append(tmp, name+"_"+strconv.Itoa(component.version[name]))
			tmp = append(tmp, component.fields[name].Interface())
		}
	}
	return tmp
}

func (component *Component) SetterCallback(args ...interface{}) bool {
	length := len(args)
	for i := 0; i < length; i += 2 {
		tmp := args[i].(string)
		splits := strings.Split(tmp, "_")
		name := splits[0]
		version, _ := strconv.Atoi(splits[1])
		value := args[i+1]
		if version > component.version[name] && component.Verify(name, value) {
			field, ok := component.fields[name]
			if ok {
				field.Set(reflect.ValueOf(value))
			}
		}
	}
	return true
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
