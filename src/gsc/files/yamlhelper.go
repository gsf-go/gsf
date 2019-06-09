package files

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func ReadFromYamlFile(fileName string, out interface{}) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, out)
	if err != nil {
		panic(err)
	}
}

func WriteToYamlFile(fileName string, in interface{}) {
	data, err := yaml.Marshal(in)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		panic(err)
	}
}
