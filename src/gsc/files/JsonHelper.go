
package files

import (
	"encoding/json"
	"io/ioutil"
)

func ReadFromJsonFile(fileName string,out interface{}){
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(data), out)
	if err != nil {
		panic(err)
	}
}

func WriteToJsonFile(fileName string,in interface{}){
	data, err := json.MarshalIndent(in, "", " ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fileName,data,0644)
	if err != nil {
		panic(err)
	}
}