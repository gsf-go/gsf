package files

import (
	"io/ioutil"
)

func ReadFromTextFile(fileName string,out interface{}){
	bytes,err:= ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	*out.(*string) = string(bytes[:])
}

func WriteToTextFile(fileName string,in interface{}){
	err:= ioutil.WriteFile(fileName,[]byte(in.(string)),0644)
	if err != nil {
		panic(err)
	}
}
