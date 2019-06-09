package files

import "testing"

func TestWriteYamlFile(t *testing.T) {
	users := []*User{
		&User{
			Name: "zs",
			Age:  50,
		}, &User{
			Name: "ls",
			Age:  100,
		},
	}
	WriteToYamlFile("file.yaml", users)
}

func TestReadYamlFile(t *testing.T) {
	users := &[]*User{}
	ReadFromYamlFile("file.yaml", users)
	for _, user := range *users {
		t.Log("Name:" + user.Name + " Age:" + string(user.Age))
	}
}
