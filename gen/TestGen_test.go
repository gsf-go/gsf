package main

//go:generate ../Generator/Singleton.exe -package=main -struct=User -out=./TestUser_test.go
type User struct {
	name string
	age int
}

//go:generate ../Generator/Enum.exe -package=main -type=int -name=TestEnum -items=One,Two -out=./TestEnum_test.go
