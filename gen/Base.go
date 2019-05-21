package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"text/template"
)

type Base struct {
	OutPath *string		`out`
	Package *string 	`package`
}

func (base *Base) Create(){
	base.OutPath = flag.String("out", ".", "输出路径！")
	base.Package = flag.String("package", "main", "包名！")
}

func (base *Base) Parse(){
	flag.Parse()
}

func GetType(v interface{}) string {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func (base *Base) Execute(name string,data interface{}) {

	files, _ := template.ParseFiles(
		filepath.Dir(os.Args[0]) + "./" + name+".tmpl")
	tmp := template.Must(files,nil)

	file, _ := os.OpenFile(*base.OutPath,
		os.O_TRUNC | os.O_CREATE | os.O_WRONLY,
		0777)

	defer file.Close()

	err := tmp.Execute(file, data)
	if err!=nil {
		fmt.Println(err)
	}
}
