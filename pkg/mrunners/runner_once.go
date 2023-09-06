package mrunners

import (
	"context"
	"github.com/guoyk93/minit/pkg/munit"
)

func init() {
	Register(munit.KindOnce, func(opts RunnerOptions) (runner Runner, err error) {
		if err = opts.Unit.RequireCommand(); err != nil {
			return
		}

		runner.Order = 20
		runner.Action = &runnerOnce{RunnerOptions: opts}
		return
	})
}

type runnerOnce struct {
	RunnerOptions
}

func (r *runnerOnce) Do(ctx context.Context) {
	r.Print("controller started")
	defer r.Print("controller exited")

	if err := r.Exec.Execute(r.Unit.ExecuteOptions(r.Logger)); err != nil {
		r.Error("failed executing: " + err.Error())
		return
	}
}
