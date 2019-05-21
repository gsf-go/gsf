package logging

import "testing"

func TestLog(t *testing.T) {
	log:= New()
	log.Info("%s", "info")
	log.Debug("%s", "debug")
	log.Warning("%s", "warning")
	log.Error("%s","error")
	log.Critical("%s","critical")
}
