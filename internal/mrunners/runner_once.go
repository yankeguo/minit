package mrunners

import (
	"context"

	"github.com/yankeguo/minit/internal/munit"
)

func init() {
	Register(munit.KindOnce, func(opts RunnerOptions) (runner Runner, err error) {
		if err = opts.Unit.RequireCommand(); err != nil {
			return
		}

		runner.Order = 20
		runner.Action = &actionOnce{RunnerOptions: opts}
		return
	})
}

type actionOnce struct {
	RunnerOptions
}

func (r *actionOnce) Do(ctx context.Context) (err error) {
	r.Print("controller started")
	defer r.Print("controller exited")

	if r.Unit.Blocking != nil && !*r.Unit.Blocking {
		go func() {
			if err := r.Execute(); err != nil {
				r.Error("failed executing (non-blocking): " + err.Error())
				return
			}
		}()
		return
	}

	if err = r.Execute(); err != nil {
		r.Error("failed executing: " + err.Error())
		if r.Unit.Critical {
			return
		} else {
			err = nil
		}
	}

	return
}
