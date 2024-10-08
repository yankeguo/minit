package mrunners

import (
	"context"

	"github.com/yankeguo/minit/internal/munit"
	"github.com/yankeguo/rg"
)

func init() {
	Register(munit.KindOnce, func(opts RunnerOptions) (runner Runner, err error) {
		defer rg.Guard(&err)
		rg.Must0(opts.Unit.RequireCommand())

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
	defer rg.Guard(&err)

	if r.Unit.Blocking != nil && !*r.Unit.Blocking {
		go func() {
			var err error
			defer rg.Guard(&err)
			err = r.PanicOnCritical("failed executing (non-blocking)", r.Execute())
		}()
		return
	}

	err = r.PanicOnCritical("failed executing", r.Execute())

	return
}
