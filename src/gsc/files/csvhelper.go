package files

import "github.com/gocarina/gocsv"

func ReadFromCsvFile(fileName string, out interface{}) {
	var text string
	ReadFromTextFile(fileName, &text)

	err := gocsv.UnmarshalString(text, out)
	if err != nil {
		panic(err)
	}
}

func WriteToCsvFile(fileName string, in interface{}) {
	data, err := gocsv.MarshalString(in)
	if err != nil {
		panic(err)
	}

	WriteToTextFile(fileName, data)
}
