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
	r.Print("started")
	defer r.Print("exited")

forLoop:
	for {
		if ctx.Err() != nil {
			break forLoop
		}

		var err error
		if err = r.Exec.Execute(r.Unit.ExecuteOptions(r.Logger)); err != nil {
			r.Error("failed executing:" + err.Error())
		}

		if ctx.Err() != nil {
			break forLoop
		}

		r.Print("restarting")

		timer := time.NewTimer(time.Second * 5)
		select {
		case <-timer.C:
		case <-ctx.Done():
			break forLoop
		}
	}

}
