package mrunners

import (
	"context"
	"github.com/guoyk93/minit/pkg/munit"
	"time"
)

func init() {
	Register(munit.KindDaemon, func(opts RunnerOptions) (runner Runner, err error) {
		if err = opts.Unit.RequireCommand(); err != nil {
			return
		}

		runner.Order = 40
		runner.Long = true
		runner.Func = &runnerDaemon{RunnerOptions: opts}
		return
	})
}

type runnerDaemon struct {
	RunnerOptions
}

func (r *runnerDaemon) Do(ctx context.Context) {
	r.Logger.Printf("runner started")
	defer r.Logger.Printf("runner exited")

forLoop:
	for {
		if ctx.Err() != nil {
			break forLoop
		}

		var err error
		if err = r.Exec.Execute(r.Unit.ExecuteOptions(r.Logger)); err != nil {
			r.Logger.Errorf("failed execute: %s", err.Error())
		}

		if ctx.Err() != nil {
			break forLoop
		}

		r.Logger.Printf("restart in 5s")

		timer := time.NewTimer(time.Second * 5)
		select {
		case <-timer.C:
		case <-ctx.Done():
			break forLoop
		}
	}

}
