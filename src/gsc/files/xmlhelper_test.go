package files

import "testing"

type User2 struct {
	Name string
	Age  int
}

func TestWriteXmlFile(t *testing.T) {
	users := []*User2{
		&User2{
			Name: "zs",
			Age:  50,
		}, &User2{
			Name: "ls",
			Age:  100,
		},
	}
	WriteToXmlFile("file.xml", users)
}

func TestReadXmlFile(t *testing.T) {
	users := &[]*User2{}
	ReadFromXmlFile("file.xml", users)
	for _, user := range *users {
		t.Log("Name:" + user.Name + " Age:" + string(user.Age))
	}
}
