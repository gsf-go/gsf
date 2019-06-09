package main

import (
	"flag"
	"strings"
)

type ItemArray []string

type Enum struct {
	*Base `base`
	Type  *string    `type`
	Items *ItemArray `items`
	Name  *string    `name`
}

func (itemArray *ItemArray) Set(value string) error {
	*itemArray = strings.Split(value, ",")
	return nil
}

func (itemArray *ItemArray) String() string {
	return strings.Join(*itemArray, ",")
}

func (enum *Enum) Create() {
	enum.Base.Create()

	enum.Type = flag.String("type", "int", "枚举值类型！")
	enum.Name = flag.String("name", "Enum", "枚举名称！")
	enum.Items = &ItemArray{}
	flag.Var(enum.Items, "items", "枚举值名称！")

}

func main() {
	enum := &Enum{}
	enum.Base = &Base{}
	enum.Create()
	enum.Parse()
	enum.Append(GetType(enum), enum)
}
