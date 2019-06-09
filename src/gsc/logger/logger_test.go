package logger

import "testing"

func TestLog(t *testing.T) {
	Log.Info("%s", "info")
	Log.Debug("%s", "debug")
	Log.Warning("%s", "warning")
	Log.Error("%s", "error")
	Log.Fatal("%s", "critical")
}
