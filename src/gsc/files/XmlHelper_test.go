package files

import "testing"

func TestWriteXmlFile(t *testing.T){
	users := []*User{
		&User{
			Name:"zs",
			Age:50,
		},&User{
			Name:"ls",
			Age:100,
		},
	}
	WriteToXmlFile("file.xml",users)
}

func TestReadXmlFile(t *testing.T){
	users:= &[]*User{}
	ReadFromXmlFile("file.xml",users)
	for _,user := range *users{
		t.Log("Name:"+user.Name+" Age:"+string(user.Age))
	}
}
