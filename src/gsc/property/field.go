package property

type Field struct {
	Version         int
	Name            string
	PropertyChanged func(string, interface{})

	value interface{}
}

func (field *Field) Getter() interface{} {
	return field.value
}

func (field *Field) Setter(value interface{}) {
	if field.value == value {
		return
	}
	field.value = value
	field.Version++
	if field.PropertyChanged != nil {
		field.PropertyChanged(field.Name, value)
	}
}
