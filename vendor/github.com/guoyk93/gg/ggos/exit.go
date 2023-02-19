package ggos

import (
	"github.com/guoyk93/gg"
	"os"
)

var (
	onExit = DefaultOnExit
	osExit = os.Exit
)

// ExitCoder interface that returns exit code
type ExitCoder interface {
	ExitCode() int
}

// Exit the exit method suitable for defer in main()
func Exit(err *error) {
	if onExit != nil {
		onExit(err)
	}

	if *err == nil {
		osExit(0)
		return
	}

	if ec, ok := (*err).(ExitCoder); ok {
		osExit(ec.ExitCode())
	} else {
		osExit(1)
	}
}

// OnExit change the default OnExit function
func OnExit(fn func(err *error)) {
	onExit = fn
}

// DefaultOnExit the default OnExit function, just print error message
func DefaultOnExit(err *error) {
	if *err == nil {
		return
	}
	gg.Log("exited with error:", (*err).Error())
}
