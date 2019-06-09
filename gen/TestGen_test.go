package main

//go:generate ../gen/Singleton.exe -struct=User -out=./TestGen_test.go

//go:generate ../gen/Enum.exe -package=main -type=int -name=TestEnum -items=One,Two -out=./TestEnum_test.go
