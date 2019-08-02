package logger

import (
	"testing"
	"time"
)

func TestLog(t *testing.T) {

	config := NewLogConfig()
	config.Capacity = 100

	Log.Info("%s", "info")
	Log.Debug("%s", "debug")
	Log.Warning("%s", "warning")
	Log.Error("%s", "error")
	Log.Fatal("%s", "critical")

	time.Sleep(time.Second * time.Duration(2))
	close(Log.logChan)
}
