package main

//go:generate ../gen/Singleton.exe -struct=singleton -out=./singleton_test.go

//go:generate ../gen/Enum.exe -package=main -type=int -name=enum -item=one,two -out=./enum_test.go
