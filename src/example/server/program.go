package main

import (
	"github.com/sf-go/gsf/src/gsm"
	"os"
	"time"
)

func main() {
	application := NewApplication()
	gsm.RunServer(application, os.Args)
	time.Sleep(3600 * time.Second)
}
