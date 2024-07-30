package mrunners

import (
	"context"
	"time"

	"github.com/yankeguo/minit/internal/munit"
)

func init() {
	Register(munit.KindDaemon, func(opts RunnerOptions) (runner Runner, err error) {
		if err = opts.Unit.RequireCommand(); err != nil {
			return
		}

		runner.Order = 40
		runner.Long = true
		runner.Action = &actionDaemon{RunnerOptions: opts}
		return
	})
}

type actionDaemon struct {
	RunnerOptions
}

func (r *actionDaemon) Do(ctx context.Context) (err error) {
	r.Print("controller started")
	defer r.Print("controller exited")

forLoop:
	for {
		if ctx.Err() != nil {
			break forLoop
		}

		if err = r.Execute(); err != nil {
			r.Error("failed executing:" + err.Error())

			if r.Unit.Critical {
				return
			} else {
				err = nil
			}
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

	return
}
