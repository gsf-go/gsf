package logger

import (
	"log"
	"os"
)

type (
	Logger interface {
		Debug(format string, args ...interface{})
		Info(format string, args ...interface{})
		Warning(format string, args ...interface{})
		Error(format string, args ...interface{})
		Fatal(format string, args ...interface{})
	}

	consoleLogger struct {
		debug   *log.Logger
		info    *log.Logger
		warning *log.Logger
		error   *log.Logger
		fatal   *log.Logger
	}
)

var Log Logger

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
	Log = NewConsoleLogger()

}

func NewConsoleLogger() Logger {
	return &consoleLogger{
		debug:   log.New(os.Stdout, CLR_0, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
		info:    log.New(os.Stdout, CLR_G, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
		warning: log.New(os.Stdout, CLR_Y, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
		error:   log.New(os.Stdout, CLR_R, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
		fatal:   log.New(os.Stdout, CLR_C, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
	}
}

func (logger *consoleLogger) Debug(format string, args ...interface{}) {
	logger.debug.Printf("DEBUG: "+format, args...)
}

func (logger *consoleLogger) Info(format string, args ...interface{}) {
	logger.info.Printf("INFO:  "+format, args...)
}

func (logger *consoleLogger) Warning(format string, args ...interface{}) {
	logger.warning.Printf("WARN:  "+format, args...)
}

func (logger *consoleLogger) Error(format string, args ...interface{}) {
	logger.error.Printf("ERROR: "+format, args...)
}

func (logger *consoleLogger) Fatal(format string, args ...interface{}) {
	logger.fatal.Printf("FATAL: "+format, args...)
}
