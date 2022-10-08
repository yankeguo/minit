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
		runner.Func = &runnerOnce{RunnerOptions: opts}
		return
	})
}

type runnerOnce struct {
	RunnerOptions
}

func (r *runnerOnce) Do(ctx context.Context) {
	r.Logger.Printf("runner started")
	defer r.Logger.Printf("runner exited")
	if err := r.Exec.Execute(r.Unit.ExecuteOptions(r.Logger)); err != nil {
		r.Logger.Errorf("failed executing: %s", err.Error())
		return
	}
}
