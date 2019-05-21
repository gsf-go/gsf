package files

import "testing"

type User struct {
	Name string
	Age int
}

func TestWriteCsvFile(t *testing.T){
	users := []*User{
		&User{
			Name:"zs",
			Age:50,
		},&User{
			Name:"ls",
			Age:100,
		},
	}
	WriteToCsvFile("file.csv",users)
}

func TestReadCsvFile(t *testing.T){
	users:= &[]*User{}
	ReadFromCsvFile("file.csv",users)
	for _,user := range *users{
		t.Log("Name:"+user.Name+" Age:"+string(user.Age))
	}
}
