package gg

import (
	"fmt"
	"log"
)

// Logger is a generic logging interface
type Logger interface {
	// Log inserts a log entry.  Arguments may be handled in the manner
	// of fmt.Print, but the underlying logger may also decide to handle
	// them differently.
	Log(v ...interface{})
	// Logf insets a log entry.  Arguments are handled in the manner of
	// fmt.Printf.
	Logf(format string, v ...interface{})
}

var (
	// DefaultLogger The global default logger
	DefaultLogger Logger = &sysLogger{}
)

// sysLogger is used as a placeholder for the default logger
type sysLogger struct{}

func (n *sysLogger) Log(v ...interface{}) {
	_ = log.Output(2, fmt.Sprint(v...))
}

func (n *sysLogger) Logf(format string, v ...interface{}) {
	_ = log.Output(2, fmt.Sprintf(format, v...))
}

// Log logs using the default logger
func Log(v ...interface{}) {
	DefaultLogger.Log(v...)
}

// Logf logs formatted using the default logger
func Logf(format string, v ...interface{}) {
	DefaultLogger.Logf(format, v...)
}
