package files

import "testing"

func TestWriteText(t *testing.T) {
	str := "Test"
	WriteToTextFile("file.txt", str)
}

func TestReadText(t *testing.T) {
	str := ""
	ReadFromTextFile("file.txt", &str)
	t.Log(str)
}
