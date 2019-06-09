package main

import (
	"gsm"
	"os"
)

func main() {
	application := NewApplication()
	gsm.RunServer(application, os.Args)
}
