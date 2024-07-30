package mrunners

import (
	"bytes"
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yankeguo/minit/internal/mexec"
	"github.com/yankeguo/minit/internal/mlog"
	"github.com/yankeguo/minit/internal/munit"
	"github.com/yankeguo/rg"
)

func TestRunnerDaemon(t *testing.T) {
	exem := mexec.NewManager()

	buf := &bytes.Buffer{}

	r := &actionDaemon{
		RunnerOptions: RunnerOptions{
			Unit: munit.Unit{
				Kind:  munit.KindDaemon,
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

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := r.Do(ctx)
		require.NoError(t, err)
	}()

	time.Sleep(time.Millisecond * 2500)

	ctxCancel()

	wg.Wait()

	require.Equal(t, 1, bytes.Count(buf.Bytes(), []byte("hello\n")))
}

func TestRunnerDaemonCritical(t *testing.T) {
	exem := mexec.NewManager()

	buf := &bytes.Buffer{}

	r := &actionDaemon{
		RunnerOptions: RunnerOptions{
			Unit: munit.Unit{
				Kind:  munit.KindDaemon,
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

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := r.Do(ctx)
		require.Error(t, err)
	}()

	time.Sleep(time.Millisecond * 2500)

	ctxCancel()

	wg.Wait()
}
