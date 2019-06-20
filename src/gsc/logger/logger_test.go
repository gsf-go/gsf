package logger

import "testing"

func TestLog(t *testing.T) {

	config := NewLogConfig()
	config.Capacity = 100
	Log.SetConfig(config)

	Log.Info("%s", "info")
	Log.Debug("%s", "debug")
	Log.Warning("%s", "warning")
	Log.Error("%s", "error")
	Log.Fatal("%s", "critical")

	close(Log.logChan)
}
