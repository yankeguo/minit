package mrunners

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yankeguo/minit/internal/mexec"
	"github.com/yankeguo/minit/internal/mlog"
	"github.com/yankeguo/minit/internal/munit"
	"github.com/yankeguo/rg"
)

func TestRunnerOnce(t *testing.T) {
	exem := mexec.NewManager()

	buf := &bytes.Buffer{}

	r := &runnerOnce{
		RunnerOptions: RunnerOptions{
			Unit: munit.Unit{
				Kind:  munit.KindOnce,
				Name:  "test",
				Shell: "/bin/bash",
				Command: []string{
					"sleep 1 && echo hello && exit 2",
				},
				Critical:     true,
				SuccessCodes: []int{0, 2},
			},
			Exec: exem,
			Logger: rg.Must(mlog.NewProcLogger(mlog.ProcLoggerOptions{
				ConsoleOut: buf,
				ConsoleErr: buf,
			})),
		},
	}

	err := r.Do(context.Background())
	require.NoError(t, err)
}

func TestRunnerOnceCritical(t *testing.T) {
	exem := mexec.NewManager()

	buf := &bytes.Buffer{}

	r := &runnerOnce{
		RunnerOptions: RunnerOptions{
			Unit: munit.Unit{
				Kind:  munit.KindOnce,
				Name:  "test",
				Shell: "/bin/bash",
				Command: []string{
					"sleep 1 && echo hello && exit 2",
				},
				Critical:     true,
				SuccessCodes: []int{1},
			},
			Exec: exem,
			Logger: rg.Must(mlog.NewProcLogger(mlog.ProcLoggerOptions{
				ConsoleOut: buf,
				ConsoleErr: buf,
			})),
		},
	}

	err := r.Do(context.Background())
	require.Error(t, err)
}

func TestRunnerOnceCriticalNonBlocking(t *testing.T) {
	exem := mexec.NewManager()

	buf := &bytes.Buffer{}

	blocking := false

	r := &runnerOnce{
		RunnerOptions: RunnerOptions{
			Unit: munit.Unit{
				Kind:     munit.KindOnce,
				Name:     "test",
				Shell:    "/bin/bash",
				Blocking: &blocking,
				Command: []string{
					"sleep 1 && echo hello && exit 2",
				},
				Critical:     true,
				SuccessCodes: []int{1},
			},
			Exec: exem,
			Logger: rg.Must(mlog.NewProcLogger(mlog.ProcLoggerOptions{
				ConsoleOut: buf,
				ConsoleErr: buf,
			})),
		},
	}

	start := time.Now()

	err := r.Do(context.Background())
	require.NoError(t, err)
	require.True(t, time.Since(start) < time.Millisecond*100)
}
