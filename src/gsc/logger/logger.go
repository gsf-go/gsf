package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	debug   *log.Logger
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
	fatal   *log.Logger
}

var Log *Logger

const (
	CLR_0 = "\x1b[30;1m"
	CLR_R = "\x1b[31;1m"
	CLR_G = "\x1b[32;1m"
	CLR_Y = "\x1b[33;1m"
	CLR_B = "\x1b[34;1m"
	CLR_M = "\x1b[35;1m"
	CLR_C = "\x1b[36;1m"
	CLR_W = "\x1b[37;1m"
	CLR_N = "\x1b[0m"
)

func init() {
	Log = NewLogger()
}

func NewLogger() *Logger {
	logger := &Logger{
		debug:   log.New(os.Stdout, CLR_0, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
		info:    log.New(os.Stdout, CLR_G, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
		warning: log.New(os.Stdout, CLR_Y, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
		error:   log.New(os.Stdout, CLR_R, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
		fatal:   log.New(os.Stdout, CLR_C, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
	}
	return logger
}

func (logger *Logger) Debug(format string, args ...interface{}) {
	_ = logger.debug.Output(2, fmt.Sprintf("DEBUG:  "+format, args...))
}

func (logger *Logger) Info(format string, args ...interface{}) {
	_ = logger.info.Output(2, fmt.Sprintf("INFO:  "+format, args...))
}

func (logger *Logger) Warning(format string, args ...interface{}) {
	_ = logger.warning.Output(2, fmt.Sprintf("WARN:  "+format, args...))
}

func (logger *Logger) Error(format string, args ...interface{}) {
	_ = logger.error.Output(2, fmt.Sprintf("ERROR:  "+format, args...))
}

func (logger *Logger) Fatal(format string, args ...interface{}) {
	_ = logger.fatal.Output(2, fmt.Sprintf("FATAL:  "+format, args...))
}
