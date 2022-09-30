package main

import (
	"context"
	"github.com/guoyk93/grace/gracelog"
)

type RunnerLevel int

const (
	RunnerL1 RunnerLevel = iota + 1
	RunnerL2
	RunnerL3
)

type RunnerFactory struct {
	Level  RunnerLevel
	Create func(unit Unit, logger gracelog.ProcLogger) (Runner, error)
}

var (
	RunnerFactories = map[string]*RunnerFactory{
		KindRender: {
			Level: RunnerL1,
			Create: func(unit Unit, logger gracelog.ProcLogger) (Runner, error) {
				return NewRenderRunner(unit, logger)
			},
		},
		KindOnce: {
			Level: RunnerL2,
			Create: func(unit Unit, logger gracelog.ProcLogger) (Runner, error) {
				return NewOnceRunner(unit, logger)
			},
		},
		KindDaemon: {
			Level: RunnerL3,
			Create: func(unit Unit, logger gracelog.ProcLogger) (Runner, error) {
				return NewDaemonRunner(unit, logger)
			},
		},
		KindCron: {
			Level: RunnerL3,
			Create: func(unit Unit, logger gracelog.ProcLogger) (Runner, error) {
				return NewCronRunner(unit, logger)
			},
		},
	}
)

type Runner interface {
	Run(ctx context.Context)
}
