package main

import (
	"gsm"
	"os"
)

func main() {
	application := NewApplication()
	gsm.RunClient(application, os.Args)
}
