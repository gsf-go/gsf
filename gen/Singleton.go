package main

import (
	"flag"
)

type Singleton struct {
	*Base			`base`
	Struct *string	`struct`
}

func (singleton *Singleton) Create(){
	singleton.Base.Create()
	singleton.Struct = flag.String("struct", "main", "结构名！")
}

func main() {
	singleton := &Singleton{}
	singleton.Base = &Base{}
	singleton.Create()
	singleton.Parse()
	singleton.Execute(GetType(singleton),singleton)
}