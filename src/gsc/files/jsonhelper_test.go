package files

import "testing"

func TestWriteJsonFile(t *testing.T) {
	users := []*User{
		&User{
			Name: "zs",
			Age:  50,
		}, &User{
			Name: "ls",
			Age:  100,
		},
	}
	WriteToJsonFile("file.json", users)
}

func TestReadJsonFile(t *testing.T) {
	users := &[]*User{}
	ReadFromJsonFile("file.json", users)
	for _, user := range *users {
		t.Log("Name:" + user.Name + " Age:" + string(user.Age))
	}
}
