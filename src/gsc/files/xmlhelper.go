package files

import (
	"encoding/xml"
	"io/ioutil"
)

func ReadFromXmlFile(fileName string, out interface{}) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	err = xml.Unmarshal(data, out)
	if err != nil {
		panic(err)
	}
}

func WriteToXmlFile(fileName string, in interface{}) {
	data, err := xml.MarshalIndent(in, "", "    ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		panic(err)
	}
}
