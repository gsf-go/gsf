package logger

const (
	Console = iota
	File
)

type LogConfig struct {
	Capacity int
	LogType  int
}

func NewLogConfig() *LogConfig {
	return &LogConfig{}
}
