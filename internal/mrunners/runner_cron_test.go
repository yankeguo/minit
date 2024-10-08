package mrunners

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yankeguo/minit/internal/mexec"
	"github.com/yankeguo/minit/internal/mlog"
	"github.com/yankeguo/minit/internal/munit"
	"github.com/yankeguo/rg"
)

func TestRunnerCron(t *testing.T) {
	exem := mexec.NewManager()

	buf := &bytes.Buffer{}

	r := &actionCron{
		RunnerOptions: RunnerOptions{
			Unit: munit.Unit{
				Kind:      munit.KindCron,
				Name:      "test",
				Cron:      "@every 2s",
				Immediate: true,
				Command: []string{
					"echo", "hhhlll",
				},
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

	require.Equal(t, 2, strings.Count(buf.String(), "hhhlll\n"))
}

func TestRunnerCronCritical(t *testing.T) {
	exem := mexec.NewManager()

	buf := &bytes.Buffer{}

	r := &actionCron{
		RunnerOptions: RunnerOptions{
			Unit: munit.Unit{
				Kind:     munit.KindCron,
				Name:     "test",
				Cron:     "@every 1s",
				Shell:    "/bin/bash",
				Critical: true,
				Command: []string{
					"exit 2",
				},
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

	wg.Wait()
}
